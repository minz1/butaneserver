package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/coreos/butane/config"
	"github.com/coreos/butane/config/common"
)

func createIgnitionHandler(c []byte) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment; filename=ignition.ign")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprint(len(c)))
		w.Write(c)
	}
}

func createButaneHandler(c []byte) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment; filename=butane.bu")
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", fmt.Sprint(len(c)))
		w.Write(c)
	}
}

func main() {
	port := flag.Int("p", 8080, "server port")
	configPath := flag.String("file-path", "config.bu", "path of the butane config file")
	flag.Parse()

	absConfigPath, err := filepath.Abs(*configPath)
	if (err != nil) {
		panic(err.Error())
	}

	butaneBytes, err := os.ReadFile(absConfigPath)
	if (err != nil) {
		panic(err.Error())
	}

	translationOptions := &common.TranslateOptions{FilesDir: "", NoResourceAutoCompression: false, DebugPrintTranslations: false}
	translateByteOptions := &common.TranslateBytesOptions{TranslateOptions: *translationOptions, Pretty: true, Raw: true}

	ignitionBytes, _, err := config.TranslateBytes(butaneBytes, *translateByteOptions)
	if (err != nil) {
		panic(err.Error())
	}

	ignitionHandler := createIgnitionHandler(ignitionBytes)
	butaneHandler := createButaneHandler(butaneBytes)

	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ignition.ign":
			ignitionHandler(w, r)
		case "/butane.bu":
			butaneHandler(w, r)
		default:
			http.NotFound(w, r)
		}
	}

	server := http.Server {
		Addr: fmt.Sprintf(":%d", *port),
		Handler: handler,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
		ReadHeaderTimeout: 200 * time.Millisecond,
	}

	log.Fatal(server.ListenAndServe())
}
