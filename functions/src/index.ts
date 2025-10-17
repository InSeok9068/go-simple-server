import { setGlobalOptions } from "firebase-functions/v2";
import { onCustomEventPublished } from "firebase-functions/v2/eventarc";
// Cloud Audit Logs 이벤트 payload 타입 (필수는 아님, any로 둬도 됨)
import type { LogEntryData as CloudAuditLogData } from "@google/events/cloud/audit/v1/LogEntryData";
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
 *
 * ⚠️ 중요 변경점:
 * - Cloud Storage "직접" 이벤트(google.cloud.storage.object.v1.finalized)가 아닌
 *   Cloud Audit Logs 이벤트를 구독해서, Eventarc의 "경로 패턴(Path Pattern)" 필터를 사용.
 * - Path Pattern을 통해 resourceName이 /projects/_/buckets/<bucket>/objects/diary/** 인 경우만 트리거.
 */

// 배포 시점에 기본 버킷명을 가져와 Path Pattern에 사용
function getDefaultBucket(): string | undefined {
  try {
    const cfg = JSON.parse(process.env.FIREBASE_CONFIG || "{}");
    return (
      (cfg.storageBucket as string | undefined) || process.env.GCLOUD_BUCKET
    );
  } catch {
    return process.env.GCLOUD_BUCKET;
  }
}

const PREFIX = "diary/"; // 처리 대상 경로 prefix

// Audit Log의 resourceName에서 bucket, object name을 파싱
// resourceName 예: projects/_/buckets/<bucket>/objects/<url-encoded object path>
function parseResourceName(resourceName: string): {
  bucket?: string;
  object?: string;
} {
  // 양쪽에 슬래시가 있을 수 있으므로 정규화
  const parts = resourceName.replace(/^\/+|\/+$/g, "").split("/");
  // ["projects","_","buckets","<bucket>","objects","<url-encoded>..."]
  const bIndex = parts.indexOf("buckets");
  const oIndex = parts.indexOf("objects");
  if (bIndex >= 0 && bIndex + 1 < parts.length) {
    const bucket = parts[bIndex + 1];
    let object: string | undefined;
    if (oIndex >= 0 && oIndex + 1 < parts.length) {
      // objects 뒤의 모든 세그먼트를 다시 결합 (경로에 "/"가 포함될 수 있음)
      const encoded = parts.slice(oIndex + 1).join("/");
      try {
        object = decodeURIComponent(encoded);
      } catch {
        object = encoded; // 혹시 디코딩 실패 시 원문 사용
      }
    }
    return { bucket, object };
  }
  return {};
}

// === 최종 트리거 ===
// Audit Logs 이벤트 + Path Pattern으로 정확히 diary/**만 구독
const defaultBucket = getDefaultBucket();
if (!defaultBucket) {
  logger.warn(
    "[resizeOnUpload] 기본 버킷을 찾지 못했습니다. " +
      "FIREBASE_CONFIG.storageBucket 또는 GCLOUD_BUCKET 환경변수를 설정하세요. " +
      "Path Pattern 필터가 비활성화되어 런타임에서 prefix 수동 체크로 대체합니다."
  );
}

export const resizeOnUpload = onCustomEventPublished<CloudAuditLogData | any>(
  {
    // Cloud Audit Logs
    eventType: "google.cloud.audit.log.v1.written",
    region: "us-west1",
    retry: false,

    // 정확 매칭 필터 (스토리지 오브젝트 생성에 해당)
    // storage.objects.create 가 object finalize와 동일 취지의 "생성" 시점을 커버
    eventFilters: {
      serviceName: "storage.googleapis.com",
      methodName: "storage.objects.create",
    },

    // 경로 패턴 필터 (resourceName 기준) : /projects/_/buckets/<bucket>/objects/diary/**
    // 기본 버킷을 알 수 없는 경우(로컬 등)는 생략하고, 런타임에서 다시 한 번 PREFIX 체크
    ...(defaultBucket
      ? {
          eventFilterPathPatterns: {
            resourceName: `/projects/_/buckets/${defaultBucket}/objects/${PREFIX}**`,
          } as Record<string, string>,
        }
      : {}),
  },
  async (event) => {
    // Audit Logs payload에서 resourceName을 얻는다.
    // v2 이벤트의 data.protoPayload.resourceName 또는 data.resourceName 중 하나에 존재
    const resourceName: string | undefined =
      event.data?.protoPayload?.resourceName ||
      (event.data?.resourceName as string | undefined);

    if (!resourceName) {
      logger.warn("Audit log event에 resourceName이 없습니다. 스킵.");
      return;
    }

    const { bucket: bucketName, object: name } =
      parseResourceName(resourceName);

    // 이름/타입 검증
    if (!bucketName || !name) return;

    // Path Pattern을 못 쓴 환경(기본 버킷 감지 실패 등) 대비: 런타임에서도 한 번 더 prefix 체크
    if (!name.startsWith(PREFIX)) {
      return;
    }
    if (!IMAGE_EXT.test(name)) {
      return;
    }

    const bucket = getStorage().bucket(bucketName);
    const file = bucket.file(name);

    try {
      // 메타 조회
      const [meta] = await file.getMetadata();

      // 이미 처리된 파일은 스킵 (무한루프 방지)
      if (meta?.metadata?.resized === "true") {
        logger.debug(`Skip already resized: ${name}`);
        return;
      }

      const { fmt, contentType: targetContentType } = detectFormat({
        contentType: meta.contentType,
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

      // 동일 경로 overwrite 저장 (다음 세대 finalize/생성 로그가 한 번 더 발생하나, resized=true로 즉시 스킵됨)
      await file.save(out, {
        resumable: false,
        validation: "crc32c",
        metadata: {
          contentType: targetContentType || meta.contentType || undefined,
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
