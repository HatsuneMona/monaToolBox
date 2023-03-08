package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"monaToolBox/global"
	"net/http"
)

const NO_OPERATE_LOG_KEY = "NO_OPERATE_LOG_KEY"

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func AdminOperateLog() gin.HandlerFunc {
	return func(c *gin.Context) {

		switch c.Request.Method {
		case http.MethodPost, http.MethodDelete, http.MethodPut:
			input := make(map[string]any)

			buf := bytes.NewBufferString("")

			if _, err := io.Copy(buf, c.Request.Body); err == nil {
				_ = json.Unmarshal(buf.Bytes(), &input)
				c.Request.Body = io.NopCloser(buf)
			} else {
				input["err"] = fmt.Sprintf("ShouldBindJSON err:%v", err.Error())
			}

			blw := &CustomResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw

			c.Next() // 先执行其他handler

			if v, ok := c.Get(NO_OPERATE_LOG_KEY); !ok || v.(bool) == false {
				global.Log.Sugar().Infof("admin operate log. route[%v] status[%v] input:%v, output:%v",
					c.Request.URL, c.Writer.Status(), input, blw.body.String())
			}
		}
	}
}
