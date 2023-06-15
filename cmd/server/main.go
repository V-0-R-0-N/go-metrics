package main

import (
	"context"
	"flag"
	"github.com/V-0-R-0-N/go-metrics.git/internal/environ"
	"github.com/V-0-R-0-N/go-metrics.git/internal/filer"
	"github.com/V-0-R-0-N/go-metrics.git/internal/flags"
	"github.com/V-0-R-0-N/go-metrics.git/internal/handlers"
	"github.com/V-0-R-0-N/go-metrics.git/internal/middlware/logger"
	"github.com/V-0-R-0-N/go-metrics.git/internal/router"
	"github.com/V-0-R-0-N/go-metrics.git/internal/storage"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {

	myLogger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	if myLogger != nil {
		defer myLogger.Sync()
	}
	sugar := logger.NewSugarLogger(myLogger)

	addr := flags.NetAddress{
		Host: "localhost",
		Port: 8080,
	}
	fileRestore := flags.NewFileRestore()
	flags.Server(&addr, fileRestore)
	flag.Parse()
	if err := environ.Server(&addr, fileRestore); err != nil {
		log.Fatal(err)
	}

	filer.FilerInit(fileRestore)
	if fileRestore.File != nil {
		defer fileRestore.File.Close()
	}
	sugar.Infow(
		"Server start",
		zap.String("address: ", addr.String()),
	)

	st := storage.NewStorage()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	filer.StartRestore(ctx, st, fileRestore)
	handlerStorage := handlers.NewHandlerStorage(st)

	logMiddlware := logger.Middlware{
		Log: sugar,
	}

	routerChi := router.Router(logMiddlware, handlerStorage)

	err = http.ListenAndServe(addr.String(), routerChi)
	if err != nil {
		log.Fatal(err)
	}
}
