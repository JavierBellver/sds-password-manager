package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
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

type siteData struct {
	Login        string
	Site         string
	SiteUsername string
	SitePassword string
}

type responseBody struct {
	Ok  bool
	Msg string
}

var path = "./users.txt"
var storagePath = "./storage.txt"

func createUsersFile() {
	var _, err = os.Stat(path)

	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		chk(err)
		defer file.Close()
	}
}

func createStorageFile() {
	var _, err = os.Stat(storagePath)

	if os.IsNotExist(err) {
		var file, err = os.Create(storagePath)
		chk(err)
		defer file.Close()
	}
}

func writeUser(login string, password string) {
	var file, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0660)
	chk(err)
	defer file.Close()

	_, err = file.WriteString("[login:" + login + "|")
	chk(err)
	_, err = file.WriteString("password:" + password + "]\n")
	chk(err)

	err = file.Sync()
	chk(err)
}

func writeSiteData(data siteData) {
	var file, err = os.OpenFile(storagePath, os.O_RDWR|os.O_APPEND, 0660)
	chk(err)
	defer file.Close()

	_, err = file.WriteString("[login:" + data.Login + "|")
	chk(err)
	_, err = file.WriteString("site:" + data.Site + "|")
	chk(err)
	_, err = file.WriteString("siteUsername:" + data.SiteUsername + "|")
	chk(err)
	_, err = file.WriteString("sitePassword:" + data.SitePassword + "]\n")
	chk(err)

	err = file.Sync()
	chk(err)
}

func validateUser(w http.ResponseWriter, login string, pass string) {
	file, err := os.Open(path)
	var res bool
	res = false
	s := "[login:" + login + "|password:" + pass + "]"
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(s)
		if s == scanner.Text() {
			res = true
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	response(w, res, "Resultado")
}

//DeleteFile borra el fichero
func deleteFile() {
	var err = os.Remove(path)
	chk(err)
}

func response(w io.Writer, ok bool, msg string) {
	r := responseBody{Ok: ok, Msg: msg}
	rJSON, err := json.Marshal(&r)
	chk(err)
	w.Write(rJSON)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "text/plain")
	validateUser(w, r.Form.Get("login"), r.Form.Get("password"))
}

func registroHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "text/plain")

	login := r.Form.Get("login")
	password := r.Form.Get("password")

	writeUser(login, password)
	response(w, true, "UsuarioRegistrado")
}

func storePasswordHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "text/plain")

	login := r.Form.Get("login")
	site := r.Form.Get("site")
	siteUsername := r.Form.Get("siteUsername")
	sitePassword := r.Form.Get("sitePassword")

	data := siteData{Login: login, Site: site, SiteUsername: siteUsername, SitePassword: sitePassword}

	writeSiteData(data)
	response(w, true, "Información guardada")
}

func main() {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	createUsersFile()
	createStorageFile()

	httpsMux := http.NewServeMux()

	httpsMux.Handle("/", http.HandlerFunc(homeHandler))
	httpsMux.Handle("/registro", http.HandlerFunc(registroHandler))
	httpsMux.Handle("/login", http.HandlerFunc(loginHandler))
	httpsMux.Handle("/guardarContraseña", http.HandlerFunc(storePasswordHandler))

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
