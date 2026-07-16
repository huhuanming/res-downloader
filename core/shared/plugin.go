package shared

import (
	"github.com/elazarl/goproxy"
	"net/http"
)

type Bridge struct {
	GetVersion     func() string
	GetResType     func(key string) (bool, bool)
	TypeSuffix     func(mime string) (string, string)
	MediaIsMarked  func(key string) bool
	MarkMedia      func(key string)
	RegisterSource func(mediaURL string, source SourceInfo)
	FindSource     func(mediaURL string) (SourceInfo, bool)
	GetConfig      func(key string) interface{}
	Send           func(t string, data interface{})
}

type Plugin interface {
	SetBridge(*Bridge)
	Domains() []string
	OnRequest(*http.Request, *goproxy.ProxyCtx) (*http.Request, *http.Response)
	OnResponse(*http.Response, *goproxy.ProxyCtx) *http.Response
}
