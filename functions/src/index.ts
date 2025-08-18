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
const CONCURRENCY = 5; // 동시 처리 제한
const IMAGE_EXT = /\.(jpe?g|png|webp)$/i;

// '오늘'만 처리할지, '어제'를 처리할지 선택 (스케줄 03:00 기준)
// 일상적으로는 'yesterday'가 하루치 완결 처리를 보장함.
// 요청대로 'today'로 둠.
// 'yesterday' 변경 하루전날의 데이터를 일괄처리
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
    // 어제 03:00에 전날 하루치 완결
    const y = new Date(now.getTime() - 24 * 60 * 60 * 1000);
    return `diary/${ymdKST(y)}/`;
  }
  return `diary/${ymdKST(now)}/`;
}

// getFiles의 튜플 타입 간단 정의(불필요한 외부 타입 의존 제거)
type GetFilesTuple = [
  any[],
  { pageToken?: string } | undefined,
  { nextPageToken?: string } | undefined
];

export const resizeDaily = onSchedule(
  {
    schedule: "0 3 * * *", // 매일 새벽 3시
    timeZone: "Asia/Seoul",
  },
  async () => {
    const bucket = getStorage().bucket();
    const limit = pLimit(CONCURRENCY);

    let pageToken: string | undefined = undefined;
    let processed = 0;
    let skipped = 0;

    const prefix = targetPrefix(); // ✅ 업로드일자 기준 폴더만 스캔
    logger.info(`Start resizeDaily. prefix=${prefix}`);

    do {
      const [files, nextQuery, apiResp] = (await bucket.getFiles({
        prefix, // 오늘(또는 어제) 업로드분만 대상
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

                // 다운로드 → 리사이즈(JPEG) → 덮어쓰기
                const [buf] = await file.download();
                const out = await sharp(buf)
                  .rotate()
                  .resize({ width: MAX_WIDTH, withoutEnlargement: true })
                  .jpeg({ quality: JPEG_QUALITY })
                  .toBuffer();

                await file.save(out, {
                  contentType: "image/jpeg",
                  resumable: false,
                  metadata: {
                    cacheControl: "public, max-age=31536000",
                    metadata: { resized: "true" }, // 중복 방지 마커
                  },
                });

                processed++;
                logger.log("Resized:", file.name);
              } catch (e: any) {
                logger.error(`Resize failed for ${file.name}`, e?.message || e);
              }
            })
          )
      );

      // 필요 시 처리 상한선 두기 (타임아웃 회피)
      // if (processed >= 2000) break;
    } while (pageToken);

    logger.info(
      `Done. processed=${processed}, skipped=${skipped}, day=${PROCESS_DAY}, prefix=${prefix}`
    );
  }
);
