package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/commercionetwork/didcomauth"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"

	"github.com/commercionetwork/dsb/environment"
	"github.com/natefinch/lumberjack"
)

func main() {
	variables, err := environment.Get()
	if err != nil {
		log.Fatal(err)
	}

	createStorageDir(variables.StoragePath)

	lumb := setupLogging(variables)
	mw := io.MultiWriter(os.Stderr, lumb)
	log.SetOutput(mw)

	rm := newResourceManager(variables.RedisAddr, variables.StoragePath)

	r := mux.NewRouter()

	r.HandleFunc("/get/{id:(?:.+)}", rm.HandleGetDocument)

	paths := []didcomauth.ProtectedMapping{
		{
			Methods: []string{http.MethodPost},
			Path:    "/upload/{id:(?:.+)}",
			Handler: rm.HandleUpload,
		},
	}

	err = didcomauth.Configure(
		didcomauth.Config{
			RedisHost:      variables.RedisAddr,
			ProtectedPaths: paths,
			JWTSecret:      variables.JWTSecret,
			CommercioLCD:   variables.CommercioLCD,
			CacheType:      didcomauth.CacheType(variables.CacheType),
		},
		r,
	)

	if err != nil {
		log.Fatal(err)
	}

	handlersChain := handlers.CombinedLoggingHandler(mw, r)
	handlersChain = handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(handlersChain)
	handlersChain = handlers.CompressHandler(handlersChain)

	hs := &http.Server{
		Handler:      handlersChain,
		Addr:         variables.ListenAddr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	err = hs.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func setupLogging(v environment.Variables) io.Writer {
	lumb := &lumberjack.Logger{
		Filename:   v.LogPath,
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	if v.Debug {
		log.Println("debugging enabled!")
	}

	return lumb
}

func createStorageDir(path string) {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) || !stat.IsDir() {
		if err := os.MkdirAll(path, 0755); err != nil {
			panic(fmt.Sprintf("could not create storage directory: %s", err.Error()))
		}
	}
}
