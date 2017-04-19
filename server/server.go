package main

import (
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

// See alternate IV creation from ciphertext below
//var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

// encrypt string to base64 crypto using AES
func encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) string {
	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	chk(err)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
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
	key := []byte("example key 1234")
	var file, err = os.OpenFile(storagePath, os.O_RDWR|os.O_APPEND, 0660)
	chk(err)
	defer file.Close()
	usr := encrypt(key, data.Login)
	st := encrypt(key, data.Site)
	usrname := encrypt(key, data.SiteUsername)
	stpswd := encrypt(key, data.SitePassword)

	_, err = file.WriteString("[login:" + usr + "|")
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
func getPasswordHandler(w http.ResponseWriter, r *http.Request) {
	key := []byte("example key 1234")
	r.ParseForm()
	w.Header().Set("Content-Type", "text/plain")

	inFile, _ := os.Open("storage.txt")
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		result := strings.Split(scanner.Text(), "|")
		user := strings.Split(result[0], ":")
		site := strings.Split(result[1], ":")
		log := strings.Split(result[2], ":")
		pass := strings.Split(result[3], ":")[1]
		pass = strings.TrimSuffix(pass, "]")
		usr := decrypt(key, user[1])
		st := decrypt(key, site[1])
		usrname := decrypt(key, log[1])
		stpswd := decrypt(key, pass)
		fmt.Println(string(user[1]))

		if r.Form.Get("site") == string(st) && r.Form.Get("user") == string(usr) {
			result := "[login:" + string(usr) + "|" + "site:" + string(st) + "|" + "siteUsername:" + string(usrname) + "|" + "sitePassword:" + string(stpswd) + "]"
			response(w, true, string(result))
		}
	}
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
