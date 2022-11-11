package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinLog middleware for Gin
func GinLog(logWithBody bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		var requestBody string
		if logWithBody {
			if requestBodyBytes, err := ioutil.ReadAll(c.Request.Body); err == nil {
				requestBody = string(requestBodyBytes)
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBodyBytes))
			} else {
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}

		c.Next()

		logMap := map[string]interface{}{
			"uri":        c.Request.URL,
			"method":     c.Request.Method,
			"body":       requestBody,
			"status":     c.Writer.Status(),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"refer":      c.Request.Referer(),
			"consume":    time.Now().Sub(startTime).Seconds(),
		}
		if err := c.Errors.Last(); err != nil {
			logMap["error"] = zap.Error(err).String
		}

		fmt.Println(logMap)
	}
}
