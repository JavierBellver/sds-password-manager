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

//Función de encriptado

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

func login(client http.Client, login string, password string) bool {

	data := url.Values{}
	data.Set("login", login)
	data.Set("password", password)
	r, err := client.PostForm("https://localhost:8081/login", data)
	decoder := json.NewDecoder(r.Body)
	var rsp resp
	err = decoder.Decode(&rsp)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	/*chk(err)
	io.Copy(os.Stdout, r.Body)
	body, err := ioutil.ReadAll(r.Body)
	fmt.Println(body)*/
	return rsp.Ok
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

func recuperarPass(client http.Client, user string) {
	var site string
	fmt.Println("Nombre del sitio:")
	fmt.Scanf("%s\n", &site)
	data := url.Values{}
	data.Set("site", site)
	data.Set("user", user)
	r, err := client.PostForm("https://localhost:8081/recuperar", data)
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
	var usuario string
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	for {
		fmt.Println("Acciones: ")
		fmt.Println("1.Login")
		fmt.Println("2.Registro")
		fmt.Println("3.Recuperar contraseña")
		fmt.Println("4.Almacenar Contraseña")
		fmt.Scanf("%s\n", &opc)

		switch opc {
		case "1":
			var user, password string
			fmt.Println("Introduce el usuario: ")
			fmt.Scanf("%s", &user)
			fmt.Println("Introduce el password: ")
			fmt.Scanf("%s", &password)
			if login(*client, user, password) {
				usuario = user
				println(usuario)
			}
		case "2":
			registerUser(*client)
		case "3":
			recuperarPass(*client, usuario)
		case "4":
			storePassword(*client)
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
