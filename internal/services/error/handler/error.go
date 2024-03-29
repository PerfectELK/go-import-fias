package handler

import (
	"github.com/PerfectELK/go-import-fias/internal/config"
	"log"
	"runtime"
)

func ErrorHandler(err error) error {
	isDebugMode := config.GetConfig("APP_DEBUG") == "true"
	if isDebugMode {
		_, filename, line, _ := runtime.Caller(1)
		log.Fatalf("[error] %s:%d %v \n", filename, line, err)
	}
	return err
}
