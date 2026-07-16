package shared

type MediaInfo struct {
	Id          string
	Url         string
	UrlSign     string
	CoverUrl    string
	Size        float64
	Domain      string
	Classify    string
	Suffix      string
	SavePath    string
	Status      string
	DecodeKey   string
	Description string
	ContentType string
	OtherData   map[string]string
}

type SourceInfo struct {
	Url         string
	Method      string
	Headers     string
	ContentType string
	StatusCode  int
}

type ApiInfo struct {
	Id                   string
	Url                  string
	UrlSign              string
	Method               string
	Domain               string
	StatusCode           int
	ContentType          string
	Size                 float64
	MediaCount           int
	MediaUrls            []string
	Body                 string
	BodyTruncated        bool
	RequestHeaders       string
	RequestBody          string
	RequestBodyTruncated bool
	CreatedAt            string
}
