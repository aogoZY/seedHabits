package middlewares

import (
	"github.com/gin-gonic/gin"
	"seedHabits/conf"
	"seedHabits/sdk/log"
	"time"
)

func LogAccess(startTime time.Time, c *gin.Context) {
	if conf.Config.Server.Env == "dev" {
		end := time.Now()
		latency := end.Sub(startTime)
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		clientIp := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		if raw != "" {
			path = path + "?" + raw
		}
		if statusCode >= 400 {
			log.ApiLogger.Errorf("[%s] %3d| %4v | %9s |%s", method, statusCode, latency, clientIp, path)
		} else {
			log.ApiLogger.Infof("[%s] %3d| %4v | %9s |%s", method, statusCode, latency, clientIp, path)

		}

		if len(c.Errors) > 0 {
			for _, item := range c.Errors {
				log.Logger.Error(item)
			}
		}
	}
}
