package core

import (
	"testing"

	"res-downloader/core/shared"
)

func TestResourceFindSourceNormalizesDefaultHTTPSPort(t *testing.T) {
	r := &Resource{}
	source := shared.SourceInfo{Url: "https://example.test/api/course/detail"}

	r.registerSource(`https:\/\/maxen-fe.oss-cn-chengdu.aliyuncs.com:443\/maxen-teacher\/course_docs\/pro\/1732431170768_Y-27468.mp4`, source)

	got, ok := r.findSource("https://maxen-fe.oss-cn-chengdu.aliyuncs.com/maxen-teacher/course_docs/pro/1732431170768_Y-27468.mp4")
	if !ok {
		t.Fatal("expected source to be found")
	}

	if got.Url != source.Url {
		t.Fatalf("expected source url %q, got %q", source.Url, got.Url)
	}
}
