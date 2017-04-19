package main

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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
	file, err := os.Open("d:/gocode/src/sds-password-manager/server/users.txt")
	var res bool
	res = false
	s := "[login:" + login + "|password:" + pass + "]"
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if s == scanner.Text() {
			res = true
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	token := generateToken(login)
	response(w, res, token)
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
	e := r.ParseForm()
	chk(e)
	w.Header().Set("Content-Type", "text/plain")

	login := r.Form.Get("login")
	site := r.Form.Get("site")
	siteUsername := r.Form.Get("siteUsername")
	sitePassword := r.Form.Get("sitePassword")

	data := siteData{Login: login, Site: site, SiteUsername: siteUsername, SitePassword: sitePassword}
	writeSiteData(data)
	response(w, true, "Información guardada")
}
func getPasswordHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "text/plain")

	inFile, _ := os.Open("d:/gocode/src/sds-password-manager/server/storage.txt")
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		result := strings.Split(scanner.Text(), "|")
		user := strings.Split(result[0], ":")
		site := strings.Split(result[1], ":")
		//login := strings.Split(result[2], ":")
		//password := strings.Split(result[3], ":")
		if r.Form.Get("site") == site[1] && r.Form.Get("user") == user[1] {
			response(w, true, scanner.Text())
		}
	}
}

func main() {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	createUsersFile()
	createStorageFile()

	httpsMux := http.NewServeMux()

	httpsMux.Handle("/registro", http.HandlerFunc(registroHandler))
	httpsMux.Handle("/login", http.HandlerFunc(loginHandler))
	httpsMux.Handle("/guardarContraseña", validateToken(http.HandlerFunc(storePasswordHandler)))
	httpsMux.Handle("/recuperar", http.HandlerFunc(getPasswordHandler))

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
