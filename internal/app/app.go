package app

import (
	"context"
	"errors"
	"fias_to_sql/internal/config"
	"fias_to_sql/internal/services/dirs"
	"fias_to_sql/internal/services/disk"
	"fias_to_sql/internal/services/error/handler"
	"fias_to_sql/internal/services/fias"
	"fias_to_sql/internal/services/logger"
	"fias_to_sql/internal/services/shutdown"
	"fias_to_sql/internal/services/terminal"
	"fias_to_sql/migrations"
	"fias_to_sql/pkg/db"
	"os"
	"time"
)

func Run() error {
	err := dirs.InitServiceDirs()
	if err != nil {
		return handler.ErrorHandler(err)
	}
	logger.Println("begin init app")
	err = config.InitConfig()
	if err != nil {
		return handler.ErrorHandler(err)
	}

	err = terminal.ParseArgs()
	if err != nil {
		return handler.ErrorHandler(err)
	}

	usageGB, err := disk.Usage()
	if err != nil {
		return handler.ErrorHandler(err)
	}
	if usageGB.FreeGB < 70 {
		return errors.New("no space left on device")
	}
	logger.Println("init app success")

	if shutdown.CheckGracefulShutdown() {
		logger.Println("reboot after graceful shutdown")
		err := shutdown.RebootAfterGracefulShutdown()
		if err != nil {
			return handler.ErrorHandler(err)
		}
	}

	path, err := fias.GetArchivePath()
	if err != nil {
		return handler.ErrorHandler(err)
	}

	importDestination, err := fias.GetImportDestination()
	if err != nil {
		return handler.ErrorHandler(err)
	}

	if importDestination == "db" {
		logger.Println("create db and tables if not exists")
		_, err = db.GetDbInstance()
		if err != nil {
			return handler.ErrorHandler(err)
		}

		err = migrations.CreateDatabase()
		if err != nil {
			return handler.ErrorHandler(err)
		}
		err = migrations.CreateTables()
		if err != nil {
			return handler.ErrorHandler(err)
		}
		logger.Println("create db and tables success")
	}

	logger.Println("begin import")
	beginTime := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	shutdown.OnShutdown(func() {
		cancel()
		logger.Println("start shutdown")
	})

	err = fias.ImportXml(
		ctx,
		path,
		importDestination,
	)

	if ctx.Err() != nil {
		err := shutdown.MakeDump()
		if err != nil {
			return handler.ErrorHandler(err)
		}
		logger.Println("end shutdown")
		os.Exit(-1)
	}

	endTime := time.Now()
	if err != nil {
		return handler.ErrorHandler(err)
	}
	logger.Println("import success")
	logger.Println("import time ", float64(endTime.Unix()-beginTime.Unix())/60, " minutes")

	if importDestination == "db" {
		logger.Println("begin migrate data from temp to original tables")
		err = migrations.MigrateDataFromTempTables()
		if err != nil {
			return handler.ErrorHandler(err)
		}
		logger.Println("migration success")
	}

	return nil
}
