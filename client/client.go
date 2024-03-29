/*

Este programa muestra comunicarse entre cliente y servidor,
así como el uso de HTTPS (HTTP sobre TLS) mediante certificados (autofirmados).

Conceptos: JSON, TLS

*/

package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var token string
var usuario string

func chk(e error) {
	if e != nil {
		log.Println(e.Error())
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

type resp struct {
	Ok  bool
	Msg string
}

func parseResponse(body []byte) (*resp, error) {
	var r = new(resp)
	err := json.Unmarshal(body, &r)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return r, err
}

func login(client http.Client) {
	var login, password string
	fmt.Println("Introduce el usuario: ")
	fmt.Scanf("%s\n", &login)
	fmt.Println("Introduce el password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	chk(err)
	password = string(bytePassword)
	data := url.Values{}
	hashUser := hashPassword(login)
	data.Set("login", hashUser)
	hash := hashPassword(password)
	data.Set("password", hash)
	r, err := client.PostForm("https://localhost:8081/login", data)
	chk(err)
	body, err := ioutil.ReadAll(r.Body)
	chk(err)
	res, err := parseResponse([]byte(body))
	chk(err)
	if res.Ok == true {
		token = res.Msg
		usuario = login
		fmt.Println("Exito, bienvenido " + usuario)
	} else {
		fmt.Println("Error, usuario no existente")
	}
	fmt.Println()
}

func storePassword(client http.Client) {
	var site, siteUsername, sitePassword string

	fmt.Println("Introduce el nombre del sitio web: ")
	fmt.Scanf("%s\n", &site)
	fmt.Println("Introduce tu nombre de usuario del sitio web: ")
	fmt.Scanf("%s\n", &siteUsername)
	fmt.Println("Introduce la contraseña en el sitio web: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	chk(err)
	sitePassword = string(bytePassword)
	data := url.Values{}
	data.Set("login", usuario)
	data.Add("site", site)
	data.Add("siteUsername", siteUsername)
	data.Add("sitePassword", sitePassword)
	r, err := http.NewRequest("POST", "https://localhost:8081/guardarContraseña", bytes.NewBufferString(data.Encode()))
	chk(err)
	r.Header.Add("Authorization", "bearer "+token)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	res, err := client.Do(r)
	chk(err)
	io.Copy(os.Stdout, res.Body)
	fmt.Println()
}

func registerUser(client http.Client) {
	var login, password string

	fmt.Println("Introduce el usuario  ")
	fmt.Scanf("%s\n", &login)
	fmt.Println("Introduce la contraseña: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	chk(err)
	password = string(bytePassword)
	data := url.Values{}
	hashUser := hashPassword(login)
	data.Set("login", hashUser)
	pswd := hashPassword(password)
	data.Set("password", pswd)
	r, err := client.PostForm("https://localhost:8081/registro", data)
	chk(err)
	io.Copy(os.Stdout, r.Body)
	fmt.Println()
}

func recuperarPass(client http.Client, user string) {
	var site string
	fmt.Println("Nombre del sitio:")
	fmt.Scanf("%s\n", &site)
	data := url.Values{}
	data.Set("site", site)
	data.Set("user", user)
	r, err := http.NewRequest("POST", "https://localhost:8081/recuperarContraseña", bytes.NewBufferString(data.Encode()))
	chk(err)
	r.Header.Add("Authorization", "bearer "+token)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	res, err := client.Do(r)
	chk(err)
	io.Copy(os.Stdout, res.Body)
	fmt.Println()
}

func main() {
	var opc string
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	for {
		fmt.Println("Acciones: ")
		fmt.Println("1.Login")
		fmt.Println("2.Registro")

		if token != "" {
			fmt.Println("3-Guardar contraseña")
			fmt.Println("4-Recuperar contraseña")
			fmt.Println("5-Logout")
			fmt.Scanf("%s\n", &opc)
			switch opc {
			case "1":
				login(*client)
			case "2":
				registerUser(*client)
			case "3":
				storePassword(*client)
			case "4":
				recuperarPass(*client, usuario)
			case "5":
				token = ""
				fmt.Println("Logging Out")
			default:
				fmt.Println(opc)
			}
		} else {
			fmt.Scanf("%s\n", &opc)
			switch opc {
			case "1":
				login(*client)
			case "2":
				registerUser(*client)
			default:
				fmt.Println(opc)
			}
		}
	}
}
