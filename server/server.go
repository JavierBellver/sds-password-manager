package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func chk(e error) {
	if e != nil {
		panic(e)
	}
}

type responseBody struct {
	Ok  bool
	Msg string
}

func response(w io.Writer, ok bool, msg string) {
	r := responseBody{Ok: ok, Msg: msg}
	rJSON, err := json.Marshal(&r)
	chk(err)
	w.Write(rJSON)
}

func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://localhost:8082"+r.RequestURI, http.StatusMovedPermanently)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "text/plain")

	switch r.Form.Get("cmd") {
	case "hola":
		response(w, true, "Hola "+r.Form.Get("mensaje"))
	default:
		response(w, false, "Comando inválido")
	}
}

func main() {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	httpsMux := http.NewServeMux()

	httpsMux.Handle("/", http.HandlerFunc(homeHandler))

	srv := &http.Server{Addr: ":8081", Handler: httpsMux}

	go func() {
		if err := srv.ListenAndServeTLS("certificado/cert.pem", "certificado/key.pem"); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan // espera señal SIGINT
	log.Println("Apagando servidor ...")

	// apagar servidor de forma segura
	ctx, fnc := context.WithTimeout(context.Background(), 5*time.Second)
	fnc()
	srv.Shutdown(ctx)

	log.Println("Servidor detenido correctamente")
}
