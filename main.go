package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	index_controller "github.com/mohprilaksono/url-shortener/app/controllers/index-controller"
	"github.com/mohprilaksono/url-shortener/config"
	"github.com/mohprilaksono/url-shortener/utils"
)

const PORT string = ":8000"

func init() {
	file, _ := utils.LoadFile()
	defer utils.CloseFile(file)

	for {
		b := make([]byte, 1)
		written, err := file.Read(b)
		if err == io.EOF {
			break
		}

		config.NumberOfBytesWritten += int64(written)
	}
}

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", index_controller.Index)
	router.HandleFunc("POST /generate-url", index_controller.Store)
	router.HandleFunc("GET /show/{url}", index_controller.Show)
	router.HandleFunc("GET /{url}", index_controller.Go)

	server := &http.Server{
		Addr: PORT,
		Handler: router,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// gracefull shotdown
		// doing some cleanup works here

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()

	log.Println("server is running on port", PORT)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err.Error())
	}

	<-idleConnsClosed
}