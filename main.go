package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"goissue/api"
	"goissue/config"
	"goissue/models"
	"goissue/pkgs/logger"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := config.Load("config.yml"); err != nil {
		logrus.WithError(err).Panicln("载入配置错误")
	}

	if config.C.Debug {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:     true,
			FullTimestamp:   true,
			TimestampFormat: "15:04:05",
		})

		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(os.Stdout)
		logrus.AddHook(logger.NewHook())
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})

		logrus.SetLevel(logrus.InfoLevel)
	}

	ctx, cancel := context.WithCancel(context.Background())

	if err := models.Connect(ctx, config.C.DBAddr); err != nil {
		logrus.WithError(err).Panicln("连接数据库错误")
	}

	srv := http.Server{Addr: config.C.Addr, Handler: chi.ServerBaseContext(ctx, api.Router())}

	go func() {
		logrus.Println("API 服务监听地址: ", config.C.Addr)
		if err := srv.ListenAndServe(); err != nil {
			logrus.WithError(err).Panicln("API 服务监听错误")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP)
	c := <-quit
	logrus.Debugln("程序信号: ", c)

	_ = srv.Shutdown(ctx)
	cancel()

	logrus.Println("程序正常退出.")
}
