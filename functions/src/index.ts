import { setGlobalOptions } from "firebase-functions/v2";
import { onObjectFinalized } from "firebase-functions/v2/storage";
import { logger } from "firebase-functions";
import { initializeApp } from "firebase-admin/app";
import { getStorage } from "firebase-admin/storage";
import sharp from "sharp";

// === 전역 옵션 ===
setGlobalOptions({
  region: "us-west1", // 스토리지 리전과 동일
  timeoutSeconds: 60 * 1, // 1분 (최대 9분)
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
const MAX_BYTES = 600 * 1024; // 600KB 이하면 스킵
const IMAGE_EXT = /\.(jpe?g|png|webp)$/i;

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

/**
 * diary/ 폴더로 업로드된 이미지가 "최종화"되면 즉시 리사이즈/재인코딩
 * - 이미 resized=true 라벨이 있으면 스킵 (무한루프 방지)
 * - 작은 파일(<= MAX_BYTES)은 메타만 찍고 스킵
 * - 같은 객체로 overwrite 저장 → 한 번 더 finalize 되지만, 다음 트리거에서 resized=true 감지로 스킵
 */
export const resizeOnUpload = onObjectFinalized(
  {
    region: "us-west1",
    retry: false,
    eventFilters: { namePrefix: "diary/" },
    // 필요하면 특정 버킷만 필터링: eventFilters: { bucket: "your-bucket-name" },
  },
  async (event) => {
    const { name, bucket: bucketName, contentType, metadata } = event.data;

    // 이름/타입 검증
    if (!name) return;
    if (!name.startsWith("diary/")) {
      // diary/ 외 경로는 무시
      return;
    }
    if (!IMAGE_EXT.test(name)) {
      return;
    }

    // 이미 처리된 파일은 스킵
    if (metadata?.resized === "true") {
      logger.debug(`Skip already resized: ${name}`);
      return;
    }

    const bucket = getStorage().bucket(bucketName);
    const file = bucket.file(name);

    try {
      // 메타 조회
      const [meta] = await file.getMetadata();
      const { fmt, contentType: targetContentType } = detectFormat({
        contentType: meta.contentType || contentType,
        name,
      });

      if (fmt === "unsupported") {
        logger.warn(`Unsupported format: ${name}`);
        return;
      }

      // 바이트 기준 스킵: 이미 충분히 작음
      const objectSize = Number(meta.size || 0);
      if (objectSize > 0 && objectSize <= MAX_BYTES) {
        await file.setMetadata({
          metadata: { ...(meta.metadata || {}), resized: "true" },
        });
        logger.debug(`Small enough, mark resized and skip: ${name}`);
        return;
      }

      // 다운로드
      const [buf] = await file.download();
      let image = sharp(buf).rotate(); // EXIF 회전 대응

      // 기본 리사이즈 + 포맷별 인코딩
      if (fmt === "jpeg") {
        image = image
          .resize({ width: MAX_WIDTH, withoutEnlargement: true })
          .jpeg({ quality: JPEG_QUALITY });
      } else if (fmt === "png") {
        image = image
          .resize({ width: MAX_WIDTH, withoutEnlargement: true })
          .png({ compressionLevel: 9, palette: true }); // 용량 감소
      } else if (fmt === "webp") {
        image = image
          .resize({ width: MAX_WIDTH, withoutEnlargement: true })
          .webp({ quality: WEBP_QUALITY });
      }

      let out = await image.toBuffer();

      // 여전히 큰 경우: JPEG/WEBP는 품질 단계적으로 낮춤
      if (out.length > MAX_BYTES && (fmt === "jpeg" || fmt === "webp")) {
        let q = fmt === "jpeg" ? JPEG_QUALITY : WEBP_QUALITY;
        while (q - QUALITY_STEP >= MIN_QUALITY && out.length > MAX_BYTES) {
          q -= QUALITY_STEP;
          const retry = sharp(buf)
            .rotate()
            .resize({ width: MAX_WIDTH, withoutEnlargement: true });
          out =
            fmt === "jpeg"
              ? await retry.jpeg({ quality: q }).toBuffer()
              : await retry.webp({ quality: q }).toBuffer();
        }
      }
      // PNG는 위 설정으로 최대한 줄였고, 그래도 크면 포맷 유지하여 저장

      // 메타 병합 + resized 마킹
      const mergedMeta = { ...(meta.metadata || {}), resized: "true" };

      // 동일 경로 overwrite 저장 (다음 세대 finalize 한 번 더 발생하나, resized=true로 즉시 스킵됨)
      await file.save(out, {
        resumable: false,
        validation: "crc32c",
        metadata: {
          contentType: targetContentType,
          cacheControl: meta.cacheControl || "public, max-age=31536000",
          metadata: mergedMeta, // ← 사용자 정의 메타
        },
        // v7+ (@google-cloud/storage) 권장
        preconditionOpts: { ifGenerationMatch: meta.generation },
        // 만약 당신의 버전이 v6 계열이면 대신 아래 줄을 쓰세요 (윗줄은 제거)
        // ifGenerationMatch: meta.generation as any,
      });

      logger.info(`Resized: ${name}, finalBytes=${out.length}`);
    } catch (e: any) {
      logger.error(`Resize failed for ${name}`, e?.message || e);
    }
  }
);
