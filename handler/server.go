package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"seedHabits/conf"
	"seedHabits/handler/middlewares"
	"seedHabits/sdk/log"
	"time"
)

type Server struct {
	server *http.Server
}

func serverLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		middlewares.LogAccess(start, c)
	}
}

//func authorize()gin.HandlerFunc{
//	return func(c *gin.Context) {
//		if err :=servives.Authentivate(c);err !=nil{
//			log.Logger.Error(err)
//			cores.NoremalResponse(c,errno.AuthNotAllowed,err.Error())
//			c.AbortWithStatus(http.StatusOK)
//		}
//	}
//}

func (s *Server) Init() error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	setupRoute(engine)
	addr := fmt.Sprintf("%s:%d", conf.Config.Server.Host, conf.Config.Server.Port)
	s.server = &http.Server{Addr: addr, Handler: engine}
	return nil
}

func (s *Server) Launch() {
	log.Logger.Infof("server started at %s", s.server.Addr)
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Logger.Error(err)
		}
	}()
}

func (s Server) Stop() error {
	if s.server != nil {
		log.Logger.Warning("server sutdowm now...")
		timeOut := time.Duration(conf.Config.Server.ShutdownTimeout) * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeOut)
		defer cancel()
		if err := s.server.Shutdown(ctx); err != nil {
			log.Logger.Error(err)
			return err
		}
		select {
		case <-ctx.Done():
			log.Logger.Warningf("server shutdowntime of %ds", conf.Config.Server.ShutdownTimeout)

		}
		log.Logger.Warning("server existing")
	}
	return nil
}
