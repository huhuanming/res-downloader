package plugins

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/elazarl/goproxy"
)

func TestDefaultPluginExtractMediaURLsFromEscapedJSON(t *testing.T) {
	body := `{"url":"https:\/\/maxen-fe.oss-cn-chengdu.aliyuncs.com:443\/maxen-teacher\/course_docs\/pro\/1732431170768_Y-27468.mp4?x=1\u0026y=2"}`

	got := (&DefaultPlugin{}).extractMediaURLs(body)
	if len(got) != 1 {
		t.Fatalf("expected 1 media URL, got %d: %#v", len(got), got)
	}

	want := "https://maxen-fe.oss-cn-chengdu.aliyuncs.com:443/maxen-teacher/course_docs/pro/1732431170768_Y-27468.mp4?x=1&y=2"
	if got[0] != want {
		t.Fatalf("expected %q, got %q", want, got[0])
	}
}

func TestDefaultPluginRestoreVideoSnapshotURL(t *testing.T) {
	rawURL := "https://maxen-fe.oss-cn-chengdu.aliyuncs.com:443/maxen-teacher/homework/pro/1732458479075_Y-31461.mp4?x-oss-process=video/snapshot,t_0,f_jpg"

	got, contentType, ok := (&DefaultPlugin{}).restoreVideoSnapshotURL(rawURL)
	if !ok {
		t.Fatal("expected snapshot URL to be restored")
	}

	want := "https://maxen-fe.oss-cn-chengdu.aliyuncs.com:443/maxen-teacher/homework/pro/1732458479075_Y-31461.mp4"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}

	if contentType != "video/mp4" {
		t.Fatalf("expected video/mp4, got %q", contentType)
	}
}

func TestDefaultPluginCaptureRequestBodyRestoresRequestBody(t *testing.T) {
	rawBody := `{"courseId":123,"page":1}`
	req := httptest.NewRequest("POST", "https://api.example.test/course/list", strings.NewReader(rawBody))
	ctx := &goproxy.ProxyCtx{}

	gotReq, gotResp := (&DefaultPlugin{}).OnRequest(req, ctx)
	if gotResp != nil {
		t.Fatal("expected nil response")
	}

	body, err := io.ReadAll(gotReq.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != rawBody {
		t.Fatalf("expected restored request body %q, got %q", rawBody, string(body))
	}

	data, ok := ctx.UserData.(apiContextData)
	if !ok {
		t.Fatal("expected request body snapshot in proxy context")
	}
	if data.RequestBody != rawBody {
		t.Fatalf("expected captured request body %q, got %q", rawBody, data.RequestBody)
	}
	if data.RequestBodyTruncated {
		t.Fatal("expected request body not to be truncated")
	}
}
