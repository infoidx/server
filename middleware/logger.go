package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	ilogger "github.com/infoidx/logger"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func CustomLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := ilogger.WithContext(ctx.Request.Context())
		data, err := ctx.GetRawData()
		if err != nil {
			logger.WithError(err).Error("GetRawData")
		}
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		start := time.Now()
		logger.WithFields(ilogger.GenFields(map[string]interface{}{
			"method":   ctx.Request.Method,
			"uri":      ctx.FullPath(),
			"ip":       ctx.ClientIP(),
			"body":     unmarshalJSONBody(data),
			"full_uri": ctx.Request.RequestURI,
			"header":   ctx.Request.Header,
		})).Info("visit request")

		ctx.Next()

		writerWrapper := newResponseWriterWrapper(ctx.Writer)
		ctx.Writer = writerWrapper
		latency := time.Since(start).Nanoseconds()
		logger.WithFields(ilogger.GenFields(map[string]interface{}{
			"method":    ctx.Request.Method,
			"uri":       ctx.FullPath(),
			"status":    ctx.Writer.Status(),
			"size":      ctx.Writer.Size(),
			"latency":   fmt.Sprintf("%d ns", latency),
			"failed":    ctx.Errors.Errors(),
			"resp_body": unmarshalJSONBody(writerWrapper.Body()),
		})).Info("visit response")
	}
}

func newResponseWriterWrapper(writer gin.ResponseWriter) *ResponseWriterWrapper {
	return &ResponseWriterWrapper{
		ResponseWriter: writer,
		body:           bytes.NewBuffer([]byte{}),
	}
}

type ResponseWriterWrapper struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *ResponseWriterWrapper) Write(b []byte) (int, error) {
	_, err := w.body.Write(b)
	if err != nil {
		return 0, err
	}
	return w.ResponseWriter.Write(b)
}

func (w *ResponseWriterWrapper) Body() []byte {
	return w.body.Bytes()
}

func unmarshalJSONBody(body []byte) (jsonBody map[string]interface{}) {
	jsonBody = make(map[string]interface{})
	err := json.Unmarshal(body, &jsonBody)

	// 当不为 JsonBody 时
	if err != nil {
		jsonBody["_text"] = string(body)
	}

	decoder := json.NewDecoder(bytes.NewBuffer(body))
	decoder.UseNumber()
	err = decoder.Decode(&jsonBody)
	if err != nil {
		jsonBody["_error"] = err.Error()
	}
	return jsonBody
}
