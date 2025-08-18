package firebaseutil

import (
	"fmt"
	"net/url"
	"strings"
)

// parseFirebaseURL은 Firebase/GCS URL에서 bucket, object를 안정적으로 추출한다.
// 지원 형태:
//  1. gs://<bucket>/<object>
//  2. https://firebasestorage.googleapis.com/v0/b/<bucket>/o/<object>[?alt=media&token=...]
//     (또는 /v0/b/<bucket>/o?name=<object>)
//  3. https://storage.googleapis.com/download/storage/v1/b/<bucket>/o/<object>[?alt=media]
//  4. https://storage.googleapis.com/<bucket>/<object>
//  5. https://<bucket>.storage.googleapis.com/<object>
//  6. (옵션) https://storage.cloud.google.com/<bucket>/<object>
//
// nolint:cyclop
func ParseFirebaseURL(raw string) (string, string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", err
	}

	// 1) gs://bucket/object
	if u.Scheme == "gs" {
		bucket := u.Host
		objEnc := strings.TrimPrefix(u.EscapedPath(), "/")
		object, _ := url.PathUnescape(objEnc) // %2F 등 복원
		return bucket, object, nil
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return "", "", fmt.Errorf("지원하지 않는 스킴: %s", u.Scheme)
	}

	host := strings.ToLower(u.Host)
	p := strings.TrimPrefix(u.EscapedPath(), "/") // RawPath가 없을 수 있어 EscapedPath 사용
	seg := strings.Split(p, "/")
	q := u.Query()

	// 2) firebasestorage.googleapis.com/v0/b/<bucket>/o/<object>  또는 /o?name=<object>
	if host == "firebasestorage.googleapis.com" {
		if len(seg) >= 5 && seg[0] == "v0" && seg[1] == "b" && seg[3] == "o" {
			bucket := seg[2]
			// 일반 케이스: 경로에 object가 있음
			objEnc := strings.Join(seg[4:], "/") // seg[4]부터 끝까지
			if objEnc != "" {
				object, _ := url.PathUnescape(objEnc)
				return bucket, object, nil
			}
			// 드문 케이스: /o?name=<object>
			if name := q.Get("name"); name != "" {
				object, _ := url.QueryUnescape(name)
				return bucket, object, nil
			}
		}
	}

	// 3) storage.googleapis.com/download/storage/v1/b/<bucket>/o/<object>
	// 4) storage.googleapis.com/<bucket>/<object>
	if host == "storage.googleapis.com" {
		if len(seg) >= 7 && seg[0] == "download" && seg[1] == "storage" && seg[2] == "v1" && seg[3] == "b" && seg[5] == "o" {
			bucket := seg[4]
			objEnc := strings.Join(seg[6:], "/")
			object, _ := url.PathUnescape(objEnc)
			return bucket, object, nil
		}
		if len(seg) >= 2 {
			bucket := seg[0]
			objEnc := strings.Join(seg[1:], "/")
			object, _ := url.PathUnescape(objEnc)
			return bucket, object, nil
		}
	}

	// 5) <bucket>.storage.googleapis.com/<object>
	if strings.HasSuffix(host, ".storage.googleapis.com") {
		bucket := strings.TrimSuffix(host, ".storage.googleapis.com")
		objEnc := strings.TrimPrefix(u.EscapedPath(), "/")
		object, _ := url.PathUnescape(objEnc)
		return bucket, object, nil
	}

	// 6) storage.cloud.google.com/<bucket>/<object> (콘솔 뷰어 링크)
	if host == "storage.cloud.google.com" && len(seg) >= 2 {
		bucket := seg[0]
		objEnc := strings.Join(seg[1:], "/")
		object, _ := url.PathUnescape(objEnc)
		return bucket, object, nil
	}

	return "", "", fmt.Errorf("지원하지 않는 URL 형식")
}
