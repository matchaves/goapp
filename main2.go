package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

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

func Urlparams(w http.ResponseWriter, r *http.Request) {
	//pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	//param1 := -1
	key1 := r.URL.Query()["key1"]
	key2 := r.URL.Query()["key2"]
	log.Println("Url GET /param1 /key1=" + key1[0] + " key1=" + key2[0])
	if key1 == nil {
		w.Write([]byte(`{"message": "need a number on url params key1 and key2"}`))
		return
	}

	//query := r.URL.Query()
	//location := query.Get("location")

	w.Write([]byte(fmt.Sprintf(`{"param1": %s, "param2": %s }`, key1[0], key2[0])))
}

func Postparams(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		rand.Seed(time.Now().UnixNano())
		var req callApp
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)

		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		fmt.Printf("Este serviço foi chamado por: %s ID da chamada %d \n", req.Sourceapp, req.Id)
		w.Header().Set("Content-Type", "application/json")
		req.Id = rand.Intn(50)
		log.Printf("APP2 Respondendo com ID %d \n", req.Id)
		req.Type = "Response"

		req.Responseapp = "APP2 RESPONDENDO OK"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(req)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}
}

func PrintEnv(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("IP GOAPP1", os.Getenv("IP_APP1"))
		//fmt.Fprintf(w, os.Getenv("BAR"))

		fmt.Println("APP2 pronta para responder")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("APP2 com sucesso api pronta para responder - IP GOAPP1" + os.Getenv("IP_APP1"))

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}

}

func MakeRequest(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	if r.Method == "GET" {

		urlEnv := os.Getenv("IP_APP1")
		var req callApp
		id := pathParams["id"]

		intVar, _ := strconv.Atoi(id)

		url := fmt.Sprintf("http://%s/post", urlEnv)
		req.Id = intVar
		req.Type = "Request"
		req.Sourceapp = "APP2"
		jsonData, _ := json.Marshal(req)
		log.Printf("executando uma request para a variavel do configmap: %v", urlEnv)

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
		rand.Seed(time.Now().UnixNano())
		urlEnv := os.Getenv("IP_AUTH")
		url := fmt.Sprintf("http://%s/auth", urlEnv)
		var reqAuth requestAuth
		reqAuth.Id = rand.Intn(50)
		reqAuth.Requestapp = "APP2"
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

		fmt.Printf("response Status Auth %s:", respAuth.ResponseAuth)
		json.NewEncoder(w).Encode(respAuth)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}
}

func main() {
	if os.Getenv("IP_APP1") == "" {
		log.Println("Precisa ser declarado as variaveis IP_APP1 e PORT ou PORT = default 8080 !!!")
		os.Exit(1)
	}
	port := func() string {
		if os.Getenv("PORT") == "" {
			return ":8080"
		} else {
			return ":" + os.Getenv("PORT")
		}
	}()

	fmt.Println("Iniciando aplicacao APP2 na porta " + port)
	r := mux.NewRouter()
	//r.HandleFunc("/")
	//api := r.PathPrefix("/").Subrouter()
	//api.HandleFunc("", get).Methods(http.MethodGet)
	//r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "                                          APLICACAO 1")
	//})

	r.HandleFunc("/", PrintEnv)

	r.HandleFunc("/post", Postparams)

	r.HandleFunc("/param1", Urlparams)

	r.HandleFunc("/make/{id}", MakeRequest)

	r.HandleFunc("/auth", AuthReq)

	//http.ListenAndServe(":80", nil)
	log.Fatal(http.ListenAndServe(port, r))
}
