package main

import (
	"context"
	"goissue/api"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"goissue/config"
	"goissue/models"
	"goissue/pkgs/logger"

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

	e := echo.New()
	e.HideBanner = true
	e.Debug = true
	e.Renderer = &api.Template{}

	e.Use(middleware.RequestID())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	// e.Use(middleware.HTTPSRedirect())
	e.Use(middleware.BodyLimit("4M"))
	e.Use(middleware.CSRF())

	e.Static("/static", "static")

	api.InitRouter(e.Group("/"))

	go func() {
		if err := e.Start(config.C.APIAddr); err != nil {
			logrus.WithError(err).Panicln("API 服务监听错误")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	c := <-quit
	logrus.Debugln("程序信号: ", c)

	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 10*time.Second)
	defer timeoutCancel()

	_ = e.Shutdown(timeoutCtx)
	cancel()

	logrus.Println("程序正常退出.")
}
