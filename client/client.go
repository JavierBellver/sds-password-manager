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
	io.Copy(os.Stdout, r.Body)
	fmt.Println()
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
func recuperarPass(client http.Client) {
	var site string
	fmt.Println("Nombre del sitio:")
	fmt.Scanf("%s\n", &site)
	data := url.Values{}
	data.Set("site", site)
	r, err := client.PostForm("https://localhost:8081/recuperar", data)
	chk(err)
	io.Copy(os.Stdout, r.Body)
	chk(err)
	io.Copy(os.Stdout, r.Body)
	fmt.Println()

}

/***
CLIENTE
***/

// gestiona el modo cliente
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
		fmt.Println("3.Recuperar contraseña")
		fmt.Scanf("%s\n", &opc)

		switch opc {
		case "1":
			if login(*client).OK {

			}
		case "2":
			registerUser(*client)
		case "3":
			recuperarPass(*client)
		default:
			fmt.Println(opc)
		}
	}
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
}
