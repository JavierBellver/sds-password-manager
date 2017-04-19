/*

Este programa muestra comunicarse entre cliente y servidor,
así como el uso de HTTPS (HTTP sobre TLS) mediante certificados (autofirmados).

Conceptos: JSON, TLS

*/

package main

import (
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
)

var token string
var usuario string

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
	fmt.Scanf("%s\n", &password)
	data := url.Values{}
	data.Set("login", login)
	data.Set("password", password)
	r, err := client.PostForm("https://localhost:8081/login", data)
	chk(err)
	body, err := ioutil.ReadAll(r.Body)
	chk(err)
	res, err := parseResponse([]byte(body))
	chk(err)
	if res.Ok == true {
		token = res.Msg
    usuario = login
	}
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
	data.Add("site", site)
	data.Add("siteUsername", siteUsername)
	data.Add("sitePassword", sitePassword)
	log.Println(data)
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

	fmt.Println("Introduce el usuario: ")
	fmt.Scanf("%s\n", &login)
	fmt.Println("Introduce la contraseña: ")
	fmt.Scanf("%s\n", &password)
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
