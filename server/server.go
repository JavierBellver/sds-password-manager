package main

import (
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

type responseBody struct {
	Ok  bool
	Msg string
}

var path = "./users.txt"

// CreateFile crea el fichero donde guardaremos los usuarios
func createFile() {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		chk(err)
		defer file.Close()
	}
}

// WriteUser Registra a un usuario en el fichero
func writeUser(login string, password string) {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0660)
	chk(err)
	defer file.Close()

	// write some text to file
	_, err = file.WriteString("[login:" + login + "|")
	chk(err)
	_, err = file.WriteString("password:" + password + "]\n")
	chk(err)

	// save changes
	err = file.Sync()
	chk(err)
}

// ReadUser devuelve a un usuario desde el fichero
func readUser() {
	// re-open file
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	chk(err)
	defer file.Close()

	// read file
	var text = make([]byte, 1024)
	for {
		n, error := file.Read(text)
		if error != io.EOF {
			chk(error)
		}
		if n == 0 {
			break
		}
	}
	fmt.Println(string(text))
	chk(err)
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

func registroHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Content-Type", "text/plain")

	login := r.Form.Get("login")
	password := r.Form.Get("password")

	writeUser(login, password)
	response(w, true, "UsuarioRegistrado")
}

func main() {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	createFile()

	httpsMux := http.NewServeMux()

	httpsMux.Handle("/", http.HandlerFunc(homeHandler))
	httpsMux.Handle("/registro", http.HandlerFunc(registroHandler))

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
