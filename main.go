package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type callApp struct {
	Id          int
	Type        string
	Sourceapp   string
	Responseapp string
}

type requestAuth struct {
	Id         int
	Requestapp string
}

type responseAuth struct {
	Id           int
	ResponseAuth string
}

func main() {
	fmt.Printf("Iniciando aplicacao APP2")
	r := mux.NewRouter()

	r.HandleFunc("/", PrintEnv)

	//http.ListenAndServe(":80", nil)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func PrintEnv(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("FOO:", os.Getenv("URLAPP1"))
		//fmt.Fprintf(w, os.Getenv("BAR"))

		fmt.Printf("Request no / APP2 com sucesso api pronta para responder")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("GET on / APP2")

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}

}
