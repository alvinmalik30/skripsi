package middleware

import (
	"os"
	"polen/config"
	"polen/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LogRequestMiddleware(log *logrus.Logger) gin.HandlerFunc {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(file)

	startTime := time.Now()

	return func(c *gin.Context) {
		c.Next()

		duration := time.Since(startTime)
		requestLog := model.Requesting{
			StartTime:  startTime,
			Duration:   duration,
			StatusCode: c.Writer.Status(),
			ClientIp:   c.ClientIP(),
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			UserAgent:  c.Request.UserAgent(),
		}

		switch {
		case c.Writer.Status() >= 500:
			log.Error(requestLog)
		case c.Writer.Status() >= 400:
			log.Warn(requestLog)
		default:
			log.Info(requestLog)
		}
	}

}
