package middleware

import (
	"bytes"
	"fqhWeb/internal/model"
	"fqhWeb/internal/service"
	"fqhWeb/pkg/api"
	"fqhWeb/pkg/logger"
	"fqhWeb/pkg/util/auth"
	"io"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var recordService = service.Record()

var respPool = sync.Pool{
	New: func() any {
		return make([]byte, 1024)
	},
}

// var respPool sync.Pool

// func init() {
// 	respPool.New = func() any {
// 		return make([]byte, 1024)
// 	}
// }

type responseBodyWriter struct {
	// 嵌入 gin.ResponseWriter，表示它将继承 gin.ResponseWriter 的所有字段和方法
	gin.ResponseWriter
	// 用于存储响应的内容
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	// 写入到 responseBodyWriter 的 body 字段中的缓冲区
	r.body.Write(b)
	// 同时将内容存储到 responseBodyWriter 的 body 字段中的缓冲区中，以便后续获取响应内容
	return r.ResponseWriter.Write(b)
}

func UserActionRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "GET" {
			var body []byte
			cClaims, _ := c.Get("claims")
			claims, ok := cClaims.(*auth.CustomClaims)
			if !ok {
				c.JSON(401, api.Err("token携带的claims不合法", nil))
				c.Abort()
				return
			}
			var err error
			body, err = io.ReadAll(c.Request.Body)
			if err != nil {
				logger.Log().Error("UserActionRecord", "记录用户请求body失败", err)
				c.JSON(500, api.Err("记录用户请求body失败", err))
				c.Abort()
				return
			} else {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
			record := &model.ActRecord{
				Ip:       c.ClientIP(),
				Method:   c.Request.Method,
				Path:     c.Request.URL.Path,
				Agent:    c.Request.UserAgent(),
				Body:     string(body),
				UserID:   claims.User.ID,
				Username: claims.User.Username,
			}

			writer := responseBodyWriter{
				ResponseWriter: c.Writer,
				body:           &bytes.Buffer{},
			}
			c.Writer = writer

			startNow := time.Now().Local()
			c.Next()
			record.Latency = time.Since(startNow)
			record.Status = c.Writer.Status()
			// if len(c.Errors) < 1 {
			// 	record.ErrorMessage = sql.NullString{String: "", Valid: false}
			// } else {
			// 	// 只获取私有类型的错误：record.ErrorMessage = sql.NullString{String: c.Errors.ByType(gin.ErrorTypePrivate).String(), Valid: true}
			// 	record.ErrorMessage = sql.NullString{String: c.Errors.String(), Valid: true}
			// }

			record.Resp = writer.body.String()
			if err = recordService.RecordCreate(record); err != nil {
				logger.Log().Error("UserActionRecord", "记录失败", err)
				c.JSON(500, api.Err("记录失败", err))
				c.Abort()
				return
			}
			if len(record.Resp) > 1024 {
				// 截断
				newBody := respPool.Get().([]byte)
				copy(newBody, record.Resp)
				record.Resp = string(newBody)
				defer respPool.Put(newBody)
			}

		} else {
			c.Next()
		}

	}
}
