package http_log_proxy

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"encoding/json"
	"github.com/andybalholm/brotli"
	"github.com/hl540/http-log-proxy/storage"
	"io"
	"log"
	"net/http"
	"time"
)

// LogRecorder 日志记录器
// 用于记录 [*http.Request] 和 [*http.Response]内容
// 在记录 [*http.Response] 后持久化
type LogRecorder struct {
	StorageProvider storage.HttpLogStorage
	log             *storage.HttpLogModel
}

// NewLogRecorder 返回一个新的记录器，想要传入持久化实例
func NewLogRecorder(storageProvider storage.HttpLogStorage) *LogRecorder {
	return &LogRecorder{
		StorageProvider: storageProvider,
		log:             &storage.HttpLogModel{},
	}
}

// WriteRequest 往记录器写入 [*http.Request] URL、Method、Header、Body
func (r *LogRecorder) WriteRequest(req *http.Request) {
	r.log.RequestUrl = req.URL.String()
	r.log.RequestMethod = req.Method
	r.log.RequestHeader = HeaderMarshal(req.Header)
	if req.Body != nil {
		requestBodyBytes, _ := io.ReadAll(req.Body)
		r.log.RequestBody = string(requestBodyBytes)
		req.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
	}
}

// WriteResponse 往记录器写入 [*http.Response] StatusCode、Header、Body
func (r *LogRecorder) WriteResponse(resp *http.Response) {
	r.log.ResponseCode = resp.StatusCode
	r.log.ResponseHeader = HeaderMarshal(resp.Header)
	if resp.Body != nil {
		responseBodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			r.log.RequestBody = err.Error()
		} else {
			contentEncoding := resp.Header.Get("Content-Encoding")
			encodeResponseBodyBytes := r.decodeResponseBody(responseBodyBytes, contentEncoding)
			r.log.ResponseBody = string(encodeResponseBodyBytes)

			resp.Body = io.NopCloser(bytes.NewBuffer(responseBodyBytes))
		}
	}
}

// 解码响应内容
// 通过Content-Encoding，支持br、gzip、deflate
func (r *LogRecorder) decodeResponseBody(body []byte, contentEncoding string) []byte {
	var bodyReader io.Reader = bytes.NewReader(body)
	switch contentEncoding {
	case "br":
		bodyReader = brotli.NewReader(bodyReader)
	case "gzip":
		bodyReader, _ = gzip.NewReader(bodyReader)
	case "deflate":
		bodyReader = flate.NewReader(bodyReader)
	default:
		return body
	}
	responseBodyByte, err := io.ReadAll(bodyReader)
	if err != nil {
		return []byte(err.Error())
	}
	return responseBodyByte
}

// Flush 将日志持久化到存储
func (r *LogRecorder) Flush(ctx context.Context) {
	r.log.RequestId = ctx.Value(HttpLogProxyRequestId).(string)
	r.log.AppId = ctx.Value(AppId).(string)
	r.log.CreateAt = time.Now().Unix()

	err := r.StorageProvider.AddHttpLog(ctx, r.log)
	if err != nil {
		log.Println("Error adding http log:", err)
	}
}

// HeaderMarshal 将 [http.Header] 内容json序列化
func HeaderMarshal(header http.Header) string {
	headerMap := make(map[string]string)
	for k := range header {
		headerMap[k] = header.Get(k)
	}
	jsonByte, _ := json.Marshal(headerMap)
	return string(jsonByte)
}
