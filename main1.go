package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"math/rand"
	"github.com/gorilla/mux"
)

type callApp struct {
	Id  int64
	Type string
	Sourceapp string
	Responseapp string
}

type requestAuth struct {
	Id  int64
	Requestapp string
}

type responseAuth struct {
	Id int64
	ResponseAuth string
}

func Urlparams(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	param1 := -1
	var err error
	if val, ok := pathParams["param1"]; ok {
		param1, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	param2 := -1
	if val, ok := pathParams["param2"]; ok {
		param2, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	//query := r.URL.Query()
	//location := query.Get("location")

	w.Write([]byte(fmt.Sprintf(`{"param1": %d, "param2": %d }`, param1, param2)))
}

func Postparams(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		rnd := rand.New(rand.NewSource(99))
		var req callApp
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)

		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		fmt.Printf("Este serviço foi chamado por: %s ID da chamada %d \n", req.Sourceapp, req.Id )
		w.Header().Set("Content-Type", "application/json")
		req.Id = rnd.Int63n(50)
		req.Type = "Response"
		//req.Sourceapp = req.Sourceapp
		req.Responseapp = "APP1 RESPONDENDO OK"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(req)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}
}

func PrintEnv(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("FOO:", os.Getenv("URLAPP2"))
		//fmt.Fprintf(w, os.Getenv("BAR"))

		fmt.Printf("Request no / APP1 com sucesso api pronta para responder")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("GET on / APP1")

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}

}

func MakeRequest(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	if r.Method == "GET" {
		fmt.Printf("UrlApp1: %s", os.Getenv("URLAPP2"))
		//fmt.Println("UrlApp1:", os.Getenv("FOO"))
		//fmt.Fprintf(w, os.Getenv("BAR"))
		url := os.Getenv("URLAPP2")
		var req callApp
		id := pathParams["id"]
		var err error
		var nom int64
		if nom, _ = strconv.ParseInt(id, 16, 64); err == nil {
			fmt.Printf("%T, %v", nom, nom)
		}

		req.Id = nom
		req.Type = "Request"
		req.Sourceapp = "APP1"
		jsonData, _ := json.Marshal(req)
		fmt.Println("executando uma request para a variavel do configmap: %s", os.Getenv("URLAPP2"))
		fmt.Println("executando uma request com ID: %v", pathParams["id"])
		request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		client := &http.Client{}
		response, error := client.Do(request)
		if error != nil {
			panic(error)
		}
		defer response.Body.Close()
		var resp callApp
		json.NewDecoder(response.Body).Decode(&resp)

		fmt.Println("response Status:", response.Status)
		fmt.Println("response Headers:", response.Header)
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println("response Body:", body)
		//fmt.Printf("Got %s age %d Param %s\n", tempPlayer.Name, tempPlayer.Age, val)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}

}

func AuthReq(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		rnd := rand.New(rand.NewSource(99))
		url := os.Getenv("AUTH")
		var reqAuth requestAuth
		reqAuth.Id = rnd.Int63n(50)
		reqAuth.Requestapp = "APP1"
		jsonData, _ := json.Marshal(reqAuth)
		request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		request.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		fmt.Println("Executando Request de autenticacao")
		response, error := client.Do(request)
		if error != nil {
			panic(error)
		}
		defer response.Body.Close()
		var respAuth responseAuth
		json.NewDecoder(response.Body).Decode(&respAuth)

		fmt.Printf("response Status Auth %d:", respAuth.ResponseAuth)
		json.NewEncoder(w).Encode(respAuth)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}
}

func main() {
	fmt.Printf("Iniciando aplicacao APP1 na porta 8081")
	r := mux.NewRouter()
	//r.HandleFunc("/")
	//api := r.PathPrefix("/").Subrouter()
	//api.HandleFunc("", get).Methods(http.MethodGet)
	//r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "                                          APLICACAO 1")
	//})

	r.HandleFunc("/", PrintEnv)

	r.HandleFunc("/post", Postparams)

	r.HandleFunc("/param1/{param1}/param2/{param2}", Urlparams)

	r.HandleFunc("/make/{id}", MakeRequest)

	r.HandleFunc("/auth", AuthReq)

	//http.ListenAndServe(":80", nil)
	log.Fatal(http.ListenAndServe(":8081", r))
}
