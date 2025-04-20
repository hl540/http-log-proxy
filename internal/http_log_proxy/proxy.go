package http_log_proxy

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hl540/http-log-proxy/storage"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

const (
	AppId                 = "app_id"
	AppKey                = "app_key"
	AppTarget             = "app_target"
	HttpLogProxyRequestId = "X-http-log-proxy-request-id"
)

// HttpLogProxy http代理，并记录日志
type HttpLogProxy struct {
	StorageProvider storage.Provider
}

func NewHttpLogProxy(storageProvider storage.Provider) *HttpLogProxy {
	return &HttpLogProxy{
		StorageProvider: storageProvider,
	}
}

// 实现 [http.Handler]
func (s *HttpLogProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := s.parseContext(r); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logRecorder := NewLogRecorder(s.StorageProvider)

	proxyHandler := &httputil.ReverseProxy{}
	proxyHandler.Rewrite = s.RewriteFunc(logRecorder)
	proxyHandler.ModifyResponse = s.ModifyResponseFunc(logRecorder)

	proxyHandler.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("proxy error: %v", err)
		http.Error(w, "proxy error: "+err.Error(), http.StatusBadGateway)
	}

	proxyHandler.ServeHTTP(w, r)
}

// 解析上下文，得到应用相关信息
func (s *HttpLogProxy) parseContext(r *http.Request) error {
	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.SplitN(path, "/", 2)

	if len(parts) == 0 {
		return errors.New("missing application flag")
	}

	appInfo, err := s.StorageProvider.GetAppByKey(r.Context(), parts[0])
	if err != nil {
		return fmt.Errorf("invalid application flag: %s", parts[0])
	}

	requestId := uuid.NewString()
	r.Header.Set(HttpLogProxyRequestId, requestId)
	ctx := context.WithValue(r.Context(), AppId, appInfo.Id)
	ctx = context.WithValue(ctx, AppKey, appInfo.Key)
	ctx = context.WithValue(ctx, AppTarget, appInfo.Target)
	ctx = context.WithValue(ctx, HttpLogProxyRequestId, requestId)
	*r = *r.WithContext(ctx)
	return nil
}

// RewriteFunc 重写[*http.Request]
func (s *HttpLogProxy) RewriteFunc(logRecorder *LogRecorder) func(*httputil.ProxyRequest) {
	return func(req *httputil.ProxyRequest) {
		target := req.In.Context().Value(AppTarget).(string)
		appKey := req.In.Context().Value(AppKey).(string)

		parse, _ := url.Parse(target)
		req.SetXForwarded()
		req.SetURL(parse)
		req.Out.URL.Path = strings.ReplaceAll(req.Out.URL.Path, "/"+appKey, "")

		log.Printf("%s => %s", req.In.URL.String(), req.Out.URL.String())
		logRecorder.WriteRequest(req.Out)
	}
}

// ModifyResponseFunc 重写 [*http.Response]
func (s *HttpLogProxy) ModifyResponseFunc(logRecorder *LogRecorder) func(*http.Response) error {
	return func(resp *http.Response) error {
		logRecorder.WriteResponse(resp)
		logRecorder.Flush(resp.Request.Context())
		return nil
	}
}
