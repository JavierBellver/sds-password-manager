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
		log.Println(e.Error())
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

type siteData struct {
	ID           string
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
var key = generateRandomBytes(32)

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

func writeUser(login string, pswHash string, salt string) {

	var file, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0660)
	chk(err)
	defer file.Close()

	_, err = file.WriteString("[login:" + login + "|")
	chk(err)
	_, err = file.WriteString("password:" + pswHash + "|")
	chk(err)
	_, err = file.WriteString("salt:" + salt + "]\n")
	chk(err)

	err = file.Sync()
	chk(err)
}

func writeSiteData(data siteData) {
	var file, err = os.OpenFile(storagePath, os.O_RDWR|os.O_APPEND, 0660)
	chk(err)
	defer file.Close()
	usr := encrypt(key, data.Login)
	st := encrypt(key, data.Site)
	usrname := encrypt(key, data.SiteUsername)
	stpswd := encrypt(key, data.SitePassword)

	_, err = file.WriteString("[id:" + data.ID + "|")
	chk(err)
	_, err = file.WriteString("login:" + usr + "|")
	chk(err)
	_, err = file.WriteString("site:" + st + "|")
	chk(err)
	_, err = file.WriteString("siteUsername:" + usrname + "|")
	chk(err)
	_, err = file.WriteString("sitePassword:" + stpswd + "]\n")
	chk(err)

	err = file.Sync()
	chk(err)
}

func validateUser(w http.ResponseWriter, login string, pswd string) {
	file, err := os.Open("users.txt")
	chk(err)
	var res = false
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Split(bufio.ScanLines)
	for scanner.Scan() && !res {
		result := strings.Split(scanner.Text(), "|")
		if len(result) > 1 {
			login := strings.Split(result[0], ":")
			pswdHashed := strings.Split(result[1], ":")
			salt := strings.Split(result[2], ":")[1]
			salt = strings.TrimSuffix(salt, "]")
			if checkHashedPassword(pswdHashed[1], pswd, salt) {
				res = true
				token := generateToken(login[1])
				response(w, res, token)
			}
		}
	}
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
	hashed, salt := hashPassword(password)

	writeUser(login, hashed, salt)
	response(w, true, "UsuarioRegistrado")
}

func storePasswordHandler(w http.ResponseWriter, r *http.Request) {
	e := r.ParseForm()
	chk(e)
	w.Header().Set("Content-Type", "text/plain")

	id := generateRandomString(12)
	login := r.Form.Get("login")
	site := r.Form.Get("site")
	siteUsername := r.Form.Get("siteUsername")
	sitePassword := r.Form.Get("sitePassword")

	data := siteData{ID: id, Login: login, Site: site, SiteUsername: siteUsername, SitePassword: sitePassword}
	writeSiteData(data)
	response(w, true, "Información guardada")
}

func getPasswordHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "text/plain")

	inFile, _ := os.Open("storage.txt")
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		result := strings.Split(scanner.Text(), "|")
		if len(result) > 1 {
			id := strings.Split(result[0], ":")[1]
			user := strings.Split(result[1], ":")
			site := strings.Split(result[2], ":")
			log := strings.Split(result[3], ":")
			pass := strings.Split(result[4], ":")[1]
			pass = strings.TrimSuffix(pass, "]")
			usr := decrypt(key, user[1])
			st := decrypt(key, site[1])
			usrname := decrypt(key, log[1])
			stpswd := decrypt(key, pass)

			if r.Form.Get("site") == string(st) && r.Form.Get("user") == string(usr) {
				result := "[id:" + string(id) + "|" + "login:" + string(usr) + "|" + "site:" + string(st) + "|" + "siteUsername:" + string(usrname) + "|" + "sitePassword:" + string(stpswd) + "]"
				response(w, true, string(result))
			}
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
	httpsMux.Handle("/recuperarContraseña", validateToken(http.HandlerFunc(getPasswordHandler)))

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
