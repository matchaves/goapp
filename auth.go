package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"math/rand"
	"github.com/gorilla/mux"
	"time"
)

type requestAuth struct {
	Id  int
	Requestapp string
}

type responseAuth struct {
	Id int
	ResponseAuth string
}

func AuthRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("Auth solicitado por Aplicacao")
		rand.Seed(time.Now().UnixNano())
		var req requestAuth
		var resp responseAuth
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)

		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		fmt.Printf("Auth solicitado por Aplicacao: %s | ID da Req %d \n", req.Requestapp, req.Id )
		w.Header().Set("Content-Type", "application/json")
		resp.Id = rand.Intn(50)
		resp.ResponseAuth = "Auth OK"
		//req.Sourceapp = req.Sourceapp
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}
}

func main() {
	fmt.Printf("Iniciando aplicacao de Autenticação")
	r := mux.NewRouter()
	//r.HandleFunc("/")
	//api := r.PathPrefix("/").Subrouter()
	//api.HandleFunc("", get).Methods(http.MethodGet)
	//r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "                                          APLICACAO 1")
	//})

	r.HandleFunc("/auth", AuthRoute)

	//http.ListenAndServe(":80", nil)
	log.Fatal(http.ListenAndServe(":80", r))
}