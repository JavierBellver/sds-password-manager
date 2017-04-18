/*

Este programa muestra comunicarse entre cliente y servidor,
así como el uso de HTTPS (HTTP sobre TLS) mediante certificados (autofirmados).

Conceptos: JSON, TLS

*/

package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var token string

// función para comprobar errores (ahorra escritura)
func chk(e error) {
	if e != nil {
		panic(e)
	}
}

// respuesta del servidor
type resp struct {
	Ok  bool   // true -> correcto, false -> error
	Msg string // mensaje adicional
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
	fmt.Scanf("%s", &password)
	data := url.Values{}
	data.Set("login", login)
	data.Set("password", password)
	r, err := client.PostForm("https://localhost:8081/login", data)
	chk(err)
	body, err := ioutil.ReadAll(r.Body)
	chk(err)
	res, err := parseResponse([]byte(body))
	chk(err)
	token = res.Msg
	io.Copy(os.Stdout, r.Body)
	fmt.Println()
}

func storePassword(client http.Client) {
	var login, site, siteUsername, sitePassword string
	r, err := http.NewRequest("POST", "https://localhost:8081/guardarContraseña", nil)
	chk(err)

	fmt.Println("Introduce el nombre de usuario: ")
	fmt.Scanf("%s\n", &login)
	fmt.Println("Introduce el nombre del sitio web: ")
	fmt.Scanf("%s\n", &site)
	fmt.Println("Introduce tu nombre de usuario del sitio web: ")
	fmt.Scanf("%s\n", &siteUsername)
	fmt.Println("Introduce la contraseña en el sitio web: ")
	fmt.Scanf("%s\n", &sitePassword)
	data := url.Values{}
	data.Set("login", login)
	data.Set("site", site)
	data.Set("siteUsername", siteUsername)
	data.Set("sitePassword", sitePassword)

	r.PostForm = data
	r.Header.Add("Authorization", "bearer "+token)
	res, err := client.Do(r)
	chk(err)
	io.Copy(os.Stdout, res.Body)
	fmt.Println()
}

func registerUser(client http.Client) {
	var login, password string

	fmt.Println("Introduce el usuario: ")
	fmt.Scanf("%s\n", &login)
	fmt.Println("Introduce la contraseña: ")
	fmt.Scanf("%s", &password)
	data := url.Values{}
	data.Set("login", login)
	data.Set("password", password)
	r, err := client.PostForm("https://localhost:8081/registro", data)
	chk(err)
	io.Copy(os.Stdout, r.Body)
	fmt.Println()
}

// gestiona el modo cliente
func main() {

	var opc string

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	fmt.Println("Acciones: ")
	fmt.Println("1.Login")
	fmt.Println("2.Registro")
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
