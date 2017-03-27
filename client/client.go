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
	"net/http"
	"net/url"
	"os"
)

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

// función para escribir una respuesta del servidor
func response(w io.Writer, ok bool, msg string) {
	r := resp{Ok: ok, Msg: msg}    // formateamos respuesta
	rJSON, err := json.Marshal(&r) // codificamos en JSON
	chk(err)                       // comprobamos error
	w.Write(rJSON)                 // escribimos el JSON resultante
}

func storePassword(client http.Client) {
	var login, site, siteUsername, sitePassword string

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

	r, err := client.PostForm("https://localhost:8081/guardarContraseña", data)
	chk(err)
	io.Copy(os.Stdout, r.Body)
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

func main() {
	client()
}

/***
CLIENTE
***/

// gestiona el modo cliente
func client() {
	/*var text string

	fmt.Println("Introduce texto: ")
	fmt.Scanf("%s", &text)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	data := url.Values{}      // estructura para contener los valores
	data.Set("cmd", "hola")   // comando (string)
	data.Set("mensaje", text) // usuario (string)
	fmt.Println(data)

	r, err := client.PostForm("https://localhost:8081", data) // enviamos por POST
	chk(err)
	io.Copy(os.Stdout, r.Body) // mostramos el cuerpo de la respuesta (es un reader)
	fmt.Println()*/

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	registerUser(*client)
}
