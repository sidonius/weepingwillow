package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/BurntSushi/toml"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

var elog *zap.SugaredLogger
var cnf Config
var wg *sync.WaitGroup
var mdb *mongo.Database
var authdb *mongo.Database

func contextWithCancelFunc(ctx context.Context, callback func()) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(c)

		select {
		case <-ctx.Done():
		case <-c:
			cancel()
			callback()
		}
	}()

	return ctx
}

func init() {
	if _, err := toml.DecodeFile(GetTOMLConfigFilePath(), &cnf); err != nil {
		fmt.Print(err)
	}

	runMode := Production
	if !cnf.Log.ProductionMode {
		runMode = Development
	}

	if cnf.Log.OutputPath == "" {
		cnf.Log.OutputPath = GetLogFilePath()
	}
	if cnf.Log.Size == 0 {
		cnf.Log.Size = default_log_size
	}
	if cnf.Log.Backups == 0 {
		cnf.Log.Backups = default_log_backups
	}
	if cnf.Log.Age == 0 {
		cnf.Log.Age = default_log_age
	}

	InitLog(runMode, cnf.Log.OutputPath, cnf.Log.Size, cnf.Log.Backups, cnf.Log.Age, cnf.Log.ToConsole)
	elog, _ = GetLogger()

	elog.Info("---->>>>---->>>>---->>>>----")
}

func main() {

	// --------------------------------
	wg = &sync.WaitGroup{}
	finished := make(chan bool, 1)

	ctx := contextWithCancelFunc(context.Background(), func() {
		elog.Info("-x- waiting cancel() -x-")
		wg.Wait()
		close(finished)
	})

	wg.Add(1)
	go RestServiceFunc(ctx)

	elog.Info(" \xF0\x9F\x8D\xB5 ") // tea
	<-finished

	elog.Info(" \xF0\x9F\x8D\x94 ") // hamburger
}
