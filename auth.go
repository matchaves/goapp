package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"math/rand"
	"github.com/gorilla/mux"
)

type requestAuth struct {
	Id  int64
	Requestapp string
}

type responseAuth struct {
	Id int64
	ResponseAuth string
}

func AuthRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("Auth solicitado por Aplicacao")
		rnd := rand.New(rand.NewSource(99))
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
		resp.Id = rnd.Int63n(50)
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
	fmt.Printf("Iniciando aplicacao de Autenticação na porta 8083")
	r := mux.NewRouter()
	//r.HandleFunc("/")
	//api := r.PathPrefix("/").Subrouter()
	//api.HandleFunc("", get).Methods(http.MethodGet)
	//r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "                                          APLICACAO 1")
	//})

	r.HandleFunc("/auth", AuthRoute)

	//http.ListenAndServe(":80", nil)
	log.Fatal(http.ListenAndServe(":8083", r))
}