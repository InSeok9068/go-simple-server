import { setGlobalOptions } from "firebase-functions/v2";
import { onSchedule } from "firebase-functions/v2/scheduler";
import { logger } from "firebase-functions";
import { initializeApp } from "firebase-admin/app";
import { getStorage } from "firebase-admin/storage";
import sharp from "sharp";
import pLimit from "p-limit";

// === 전역 옵션 ===
setGlobalOptions({
  region: "us-west1", // 스토리지와 동일 리전
  timeoutSeconds: 60 * 10, // 10분
  memory: "512MiB",
});

// Firebase Admin 초기화
initializeApp();

// === 설정값 ===
const MAX_WIDTH = 800;
const JPEG_QUALITY = 75;
const WEBP_QUALITY = 75;
const MIN_QUALITY = 50; // 바이트 초과 시 여기까지 낮춤
const QUALITY_STEP = 5;
const MAX_BYTES = 600 * 1024; // 600KB 이하이면 "충분히 작다"로 간주
const CONCURRENCY = 5; // 동시 처리 제한
const IMAGE_EXT = /\.(jpe?g|png|webp)$/i;

const PROCESS_DAY: "today" | "yesterday" = "yesterday";

// KST 기준 YYYY-MM-DD
function ymdKST(d = new Date()): string {
  return new Intl.DateTimeFormat("en-CA", {
    timeZone: "Asia/Seoul",
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  }).format(d);
}

// 오늘/어제 접두사 생성
function targetPrefix(): string {
  const now = new Date();
  if (PROCESS_DAY === "yesterday") {
    const y = new Date(now.getTime() - 24 * 60 * 60 * 1000);
    return `diary/${ymdKST(y)}/`;
  }
  return `diary/${ymdKST(now)}/`;
}

// getFiles의 튜플 타입 간단 정의
type GetFilesTuple = [
  any[],
  { pageToken?: string } | undefined,
  { nextPageToken?: string } | undefined
];

// 원본 메타 기준으로 포맷/콘텐츠타입 결정
function detectFormat(meta: { contentType?: string; name?: string }) {
  const ct = (meta.contentType || "").toLowerCase();
  const name = meta.name || "";
  const ext = name.split(".").pop()?.toLowerCase();

  const isJPEG = ct.includes("jpeg") || ext === "jpg" || ext === "jpeg";
  const isPNG = ct.includes("png") || ext === "png";
  const isWEBP = ct.includes("webp") || ext === "webp";

  if (isJPEG) return { fmt: "jpeg" as const, contentType: "image/jpeg" };
  if (isPNG) return { fmt: "png" as const, contentType: "image/png" };
  if (isWEBP) return { fmt: "webp" as const, contentType: "image/webp" };
  return { fmt: "unsupported" as const, contentType: meta.contentType || "" };
}

export const resizeDaily = onSchedule(
  {
    schedule: "0 3 * * *",
    timeZone: "Asia/Seoul",
  },
  async () => {
    const bucket = getStorage().bucket();
    const limit = pLimit(CONCURRENCY);

    let pageToken: string | undefined = undefined;
    let processed = 0;
    let skipped = 0;
    let unsupported = 0;

    const prefix = targetPrefix();
    logger.info(`Start resizeDaily. prefix=${prefix}`);

    do {
      const [files, nextQuery, apiResp] = (await bucket.getFiles({
        prefix,
        autoPaginate: false,
        pageToken,
      })) as GetFilesTuple;

      pageToken = nextQuery?.pageToken ?? apiResp?.nextPageToken;

      await Promise.all(
        files
          .filter((f) => IMAGE_EXT.test(f.name))
          .map((file) =>
            limit(async () => {
              try {
                // 이미 처리된 파일은 스킵
                const [meta] = await file.getMetadata();
                if (meta?.metadata?.resized === "true") {
                  skipped++;
                  return;
                }

                const { fmt, contentType } = detectFormat({
                  contentType: meta.contentType,
                  name: meta.name,
                });
                if (fmt === "unsupported") {
                  unsupported++;
                  logger.warn(`Skip unsupported format: ${file.name}`);
                  return;
                }

                // ✅ 사이즈(바이트) 기준으로 스킵
                const objectSize = Number(meta.size || 0); // GCS 메타 size는 문자열
                if (objectSize > 0 && objectSize <= MAX_BYTES) {
                  await file.setMetadata({
                    metadata: { ...(meta.metadata || {}), resized: "true" },
                  });
                  skipped++;
                  return;
                }

                // 다운로드
                const [buf] = await file.download();
                let image = sharp(buf).rotate();

                // 기본 리사이즈 (너비 제한은 유지하되, 판단 기준은 바이트)
                // -> 큰 이미지면 용량 줄이는데 유리하니 그대로 수행
                if (fmt === "jpeg") {
                  image = image
                    .resize({ width: MAX_WIDTH, withoutEnlargement: true })
                    .jpeg({ quality: JPEG_QUALITY });
                } else if (fmt === "png") {
                  image = image
                    .resize({ width: MAX_WIDTH, withoutEnlargement: true })
                    .png({ compressionLevel: 9, palette: true }); // palette로 용량 더 감소
                } else if (fmt === "webp") {
                  image = image
                    .resize({ width: MAX_WIDTH, withoutEnlargement: true })
                    .webp({ quality: WEBP_QUALITY });
                }

                let out = await image.toBuffer();

                // ✅ 인코딩 결과가 여전히 MAX_BYTES 초과이면 품질 낮춰 재시도(jpeg/webp)
                if (
                  out.length > MAX_BYTES &&
                  (fmt === "jpeg" || fmt === "webp")
                ) {
                  let q = fmt === "jpeg" ? JPEG_QUALITY : WEBP_QUALITY;
                  while (
                    q - QUALITY_STEP >= MIN_QUALITY &&
                    out.length > MAX_BYTES
                  ) {
                    q -= QUALITY_STEP;
                    const retry = sharp(buf)
                      .rotate()
                      .resize({ width: MAX_WIDTH, withoutEnlargement: true });
                    const encoded =
                      fmt === "jpeg"
                        ? retry.jpeg({ quality: q })
                        : retry.webp({ quality: q });
                    out = await encoded.toBuffer();
                  }
                }
                // PNG는 위에서 palette+compressionLevel로 최대한 줄였고,
                // 그래도 초과면 그대로 저장(포맷 유지 조건)

                // 기존 메타 병합 + resized 마킹
                const merged = { ...(meta.metadata || {}), resized: "true" };

                await file.save(out, {
                  contentType,
                  resumable: false,
                  cacheControl: meta.cacheControl || "public, max-age=31536000",
                  metadata: merged,
                  validation: "crc32c",
                });

                processed++;
                logger.log("Resized:", file.name);
              } catch (e: any) {
                logger.error(`Resize failed for ${file.name}`, e?.message || e);
              }
            })
          )
      );
    } while (pageToken);

    logger.info(
      `Done. processed=${processed}, skipped=${skipped}, unsupported=${unsupported}, day=${PROCESS_DAY}, prefix=${prefix}`
    );
  }
);
