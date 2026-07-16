package plugins

import (
	"bytes"
	"encoding/json"
	"html"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/elazarl/goproxy"
	gonanoid "github.com/matoous/go-nanoid/v2"

	"res-downloader/core/shared"
)

const sourceScanLimit = 5 * 1024 * 1024
const requestBodyCaptureLimit = 5 * 1024 * 1024

var mediaURLRegex = regexp.MustCompile(`(?i)https?://[^\s"'<>\\]+?\.(?:mp4|m3u8|flv|mov|webm)(?:\?[^\s"'<>\\]*)?`)
var encodedMediaURLRegex = regexp.MustCompile(`(?i)https?%3A%2F%2F[^\s"'<>\\]+?\.(?:mp4|m3u8|flv|mov|webm)(?:%3F[^\s"'<>\\]*)?`)

type prependReadCloser struct {
	io.Reader
	io.Closer
}

type apiContextData struct {
	RequestBody          string
	RequestBodyTruncated bool
}

type DefaultPlugin struct {
	bridge *shared.Bridge
}

func (p *DefaultPlugin) SetBridge(bridge *shared.Bridge) {
	p.bridge = bridge
}

func (p *DefaultPlugin) Domains() []string {
	return []string{"default"}
}

func (p *DefaultPlugin) OnRequest(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	r.Header.Del("Accept-Encoding")
	p.captureRequestBody(r, ctx)
	return r, nil
}

func (p *DefaultPlugin) OnResponse(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	if resp == nil || resp.Request == nil || (resp.StatusCode != 200 && resp.StatusCode != 206 && resp.StatusCode != 304) {
		return resp
	}

	p.captureSourceCandidates(resp, ctx)

	contentType := resp.Header.Get("Content-Type")
	classify, suffix := p.bridge.TypeSuffix(contentType)
	if classify == "" {
		return resp
	}

	rawUrl := resp.Request.URL.String()
	isSnapshotVideo := false
	if mediaURL, mediaContentType, ok := p.restoreVideoSnapshotURL(rawUrl); ok {
		rawUrl = mediaURL
		contentType = mediaContentType
		classify = "video"
		suffix = filepath.Ext(mediaURL)
		isSnapshotVideo = true
	}

	isAll, _ := p.bridge.GetResType("all")
	isClassify, _ := p.bridge.GetResType(classify)

	if suffix == "default" {
		ext := filepath.Ext(filepath.Base(strings.Split(strings.Split(rawUrl, "?")[0], "#")[0]))
		if ext != "" {
			suffix = ext
		}
	}

	urlSign := shared.Md5(rawUrl)
	if ok := p.bridge.MediaIsMarked(urlSign); !ok && (isAll || isClassify) {
		value, _ := strconv.ParseFloat(resp.Header.Get("content-length"), 64)
		if isSnapshotVideo {
			value, contentType = p.probeMedia(resp.Request, rawUrl, contentType)
		}

		id, err := gonanoid.New()
		if err != nil {
			id = urlSign
		}
		res := shared.MediaInfo{
			Id:          id,
			Url:         rawUrl,
			UrlSign:     urlSign,
			CoverUrl:    "",
			Size:        value,
			Domain:      shared.GetTopLevelDomain(rawUrl),
			Classify:    classify,
			Suffix:      suffix,
			Status:      shared.DownloadStatusReady,
			SavePath:    "",
			DecodeKey:   "",
			OtherData:   map[string]string{},
			Description: "",
			ContentType: contentType,
		}

		// Store entire request headers as JSON
		if headers, err := json.Marshal(resp.Request.Header); err == nil {
			res.OtherData["headers"] = string(headers)
		}
		if referer := resp.Request.Header.Get("Referer"); referer != "" {
			res.OtherData["referer"] = referer
		}
		if origin := resp.Request.Header.Get("Origin"); origin != "" {
			res.OtherData["origin"] = origin
		}
		if isSnapshotVideo {
			res.OtherData["snapshot_url"] = resp.Request.URL.String()
		}
		if source, ok := p.bridge.FindSource(rawUrl); ok {
			res.OtherData["source_url"] = source.Url
			res.OtherData["source_method"] = source.Method
			res.OtherData["source_headers"] = source.Headers
			res.OtherData["source_content_type"] = source.ContentType
			res.OtherData["source_status_code"] = strconv.Itoa(source.StatusCode)
		}

		p.bridge.MarkMedia(urlSign)
		go func(res shared.MediaInfo) {
			p.bridge.Send("newResources", res)
		}(res)
	}

	return resp
}

func (p *DefaultPlugin) captureRequestBody(r *http.Request, ctx *goproxy.ProxyCtx) {
	if r == nil || r.Body == nil || r.Body == http.NoBody {
		return
	}

	bodyPrefix, err := io.ReadAll(io.LimitReader(r.Body, requestBodyCaptureLimit+1))
	r.Body = &prependReadCloser{
		Reader: io.MultiReader(bytes.NewReader(bodyPrefix), r.Body),
		Closer: r.Body,
	}
	if err != nil || len(bodyPrefix) == 0 || ctx == nil {
		return
	}

	truncated := false
	body := bodyPrefix
	if len(body) > requestBodyCaptureLimit {
		body = body[:requestBodyCaptureLimit]
		truncated = true
	}
	if r.ContentLength > int64(len(bodyPrefix)) && r.ContentLength > requestBodyCaptureLimit {
		truncated = true
	}

	ctx.UserData = apiContextData{
		RequestBody:          string(body),
		RequestBodyTruncated: truncated,
	}
}

func (p *DefaultPlugin) captureSourceCandidates(resp *http.Response, ctx *goproxy.ProxyCtx) {
	if p.bridge == nil || resp.Body == nil || resp.Request == nil {
		return
	}

	if !p.isInspectableContent(resp.Header.Get("Content-Type")) || !p.isPlainBody(resp.Header.Get("Content-Encoding")) {
		return
	}

	bodyPrefix, err := io.ReadAll(io.LimitReader(resp.Body, sourceScanLimit))
	resp.Body = &prependReadCloser{
		Reader: io.MultiReader(bytes.NewReader(bodyPrefix), resp.Body),
		Closer: resp.Body,
	}
	if err != nil || len(bodyPrefix) == 0 {
		return
	}

	candidates := p.extractMediaURLs(string(bodyPrefix))
	if p.isJSONContent(resp.Header.Get("Content-Type")) {
		p.emitAPIInfo(resp, ctx, bodyPrefix, candidates)
	}

	if len(candidates) == 0 {
		return
	}

	sourceHeaders, _ := json.Marshal(resp.Request.Header)
	source := shared.SourceInfo{
		Url:         resp.Request.URL.String(),
		Method:      resp.Request.Method,
		Headers:     string(sourceHeaders),
		ContentType: resp.Header.Get("Content-Type"),
		StatusCode:  resp.StatusCode,
	}

	for _, candidate := range candidates {
		p.bridge.RegisterSource(candidate, source)
		if mediaURL, _, ok := p.restoreVideoSnapshotURL(candidate); ok {
			p.bridge.RegisterSource(mediaURL, source)
		}
	}
}

func (p *DefaultPlugin) isInspectableContent(contentType string) bool {
	contentType = strings.ToLower(strings.Split(contentType, ";")[0])
	return strings.HasPrefix(contentType, "text/") ||
		strings.Contains(contentType, "json") ||
		strings.Contains(contentType, "javascript") ||
		strings.Contains(contentType, "xml") ||
		contentType == "application/x-www-form-urlencoded"
}

func (p *DefaultPlugin) isJSONContent(contentType string) bool {
	contentType = strings.ToLower(strings.Split(contentType, ";")[0])
	return strings.Contains(contentType, "json")
}

func (p *DefaultPlugin) isPlainBody(contentEncoding string) bool {
	contentEncoding = strings.ToLower(strings.TrimSpace(contentEncoding))
	return contentEncoding == "" || contentEncoding == "identity"
}

func (p *DefaultPlugin) extractMediaURLs(body string) []string {
	body = html.UnescapeString(body)
	body = strings.NewReplacer(
		`\/`, `/`,
		`\u0026`, `&`,
		`\u002F`, `/`,
		`\u003a`, `:`,
		`\u003A`, `:`,
		`\u003d`, `=`,
		`\u003D`, `=`,
		`\u003f`, `?`,
		`\u003F`, `?`,
	).Replace(body)

	seen := make(map[string]bool)
	var result []string
	for _, candidate := range mediaURLRegex.FindAllString(body, -1) {
		if !seen[candidate] {
			seen[candidate] = true
			result = append(result, candidate)
		}
	}

	for _, candidate := range encodedMediaURLRegex.FindAllString(body, -1) {
		decoded, err := url.QueryUnescape(candidate)
		if err != nil {
			continue
		}
		if !seen[decoded] {
			seen[decoded] = true
			result = append(result, decoded)
		}
	}

	return result
}

func (p *DefaultPlugin) emitAPIInfo(resp *http.Response, ctx *goproxy.ProxyCtx, bodyPrefix []byte, candidates []string) {
	requestURL := resp.Request.URL.String()
	urlSign := shared.Md5(resp.Request.Method + " " + requestURL)
	id, err := gonanoid.New()
	if err != nil {
		id = urlSign
	}

	size, _ := strconv.ParseFloat(resp.Header.Get("content-length"), 64)
	if size <= 0 && resp.ContentLength > 0 {
		size = float64(resp.ContentLength)
	}
	if size <= 0 {
		size = float64(len(bodyPrefix))
	}

	truncated := false
	if len(bodyPrefix) >= sourceScanLimit && (resp.ContentLength < 0 || resp.ContentLength > int64(len(bodyPrefix))) {
		truncated = true
	}

	var requestBody string
	requestBodyTruncated := false
	if ctx != nil {
		if data, ok := ctx.UserData.(apiContextData); ok {
			requestBody = data.RequestBody
			requestBodyTruncated = data.RequestBodyTruncated
		}
	}

	sourceHeaders, _ := json.Marshal(resp.Request.Header)
	api := shared.ApiInfo{
		Id:                   id,
		Url:                  requestURL,
		UrlSign:              urlSign,
		Method:               resp.Request.Method,
		Domain:               shared.GetTopLevelDomain(requestURL),
		StatusCode:           resp.StatusCode,
		ContentType:          resp.Header.Get("Content-Type"),
		Size:                 size,
		MediaCount:           len(candidates),
		MediaUrls:            p.limitStrings(candidates, 50),
		Body:                 string(bodyPrefix),
		BodyTruncated:        truncated,
		RequestHeaders:       string(sourceHeaders),
		RequestBody:          requestBody,
		RequestBodyTruncated: requestBodyTruncated,
		CreatedAt:            time.Now().Format("15:04:05"),
	}

	go func(api shared.ApiInfo) {
		p.bridge.Send("newApis", api)
	}(api)
}

func (p *DefaultPlugin) limitStrings(values []string, limit int) []string {
	if len(values) <= limit {
		return values
	}
	result := make([]string, limit)
	copy(result, values[:limit])
	return result
}

func (p *DefaultPlugin) restoreVideoSnapshotURL(rawURL string) (string, string, bool) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", false
	}

	query := parsedURL.Query()
	processValues, ok := query["x-oss-process"]
	if !ok {
		return "", "", false
	}

	hasVideoSnapshot := false
	for _, value := range processValues {
		if strings.HasPrefix(strings.ToLower(value), "video/snapshot") {
			hasVideoSnapshot = true
			break
		}
	}
	if !hasVideoSnapshot {
		return "", "", false
	}

	contentType, ok := p.videoContentTypeByPath(parsedURL.Path)
	if !ok {
		return "", "", false
	}

	query.Del("x-oss-process")
	parsedURL.RawQuery = query.Encode()
	parsedURL.Fragment = ""
	return parsedURL.String(), contentType, true
}

func (p *DefaultPlugin) videoContentTypeByPath(path string) (string, bool) {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".mp4":
		return "video/mp4", true
	case ".mov":
		return "video/quicktime", true
	case ".webm":
		return "video/webm", true
	case ".flv":
		return "video/x-flv", true
	case ".m3u8":
		return "application/vnd.apple.mpegurl", true
	default:
		return "", false
	}
}

func (p *DefaultPlugin) probeMedia(sourceReq *http.Request, mediaURL, fallbackContentType string) (float64, string) {
	req, err := http.NewRequest(http.MethodHead, mediaURL, nil)
	if err != nil {
		return 0, fallbackContentType
	}

	req.Header = sourceReq.Header.Clone()
	req.Header.Del("Range")
	req.Header.Del("Accept-Encoding")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fallbackContentType
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return 0, fallbackContentType
	}

	contentType := resp.Header.Get("Content-Type")
	if classify, _ := p.bridge.TypeSuffix(contentType); classify == "" {
		contentType = fallbackContentType
	}

	value, _ := strconv.ParseFloat(resp.Header.Get("content-length"), 64)
	return value, contentType
}
