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

	"github.com/gorilla/mux"
)

type player struct {
	Name string
	Age  uint8
	Param string
	
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
		var tempPlayer player
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&tempPlayer)

		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		fmt.Printf("Este serviço foi chamado por: %s urlParamReq %d Param %s\n", tempPlayer.Name, tempPlayer.Age, tempPlayer.Param)
		w.Header().Set("Content-Type", "application/json")
		tempPlayer.Name = "APP2 RESPONDEU"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tempPlayer)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}
}

func PrintEnv(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("FOO:", os.Getenv("FOO"))
		//fmt.Fprintf(w, os.Getenv("BAR"))
		var tempPlayer player
		tempPlayer.Name = "APP2"
		tempPlayer.Age = 2
		tempPlayer.Param = "Sucesso pelo metodo GET no /"
		fmt.Printf("Got %s age %d Param %s\n", tempPlayer.Name, tempPlayer.Age, tempPlayer.Param)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("GET on / APP2")

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}

}

func MakeRequest(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	if r.Method == "GET" {
		fmt.Printf("UrlApp1: %s", os.Getenv("FOO"))
		//fmt.Println("UrlApp2:", os.Getenv("FOO"))
		//fmt.Fprintf(w, os.Getenv("BAR"))
		url := os.Getenv("FOO")
		val := pathParams["app"]
		var tempPlayer player
		tempPlayer.Name = val
		tempPlayer.Age = 1
		tempPlayer.Param = "request do app1"
		jsonData, _ := json.Marshal(tempPlayer)
		fmt.Println("variavel recebida pelo pathParam: %s", val)
		fmt.Println("executando uma request para a variavel do configmap: %s", os.Getenv("FOO"))
		request, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		client := &http.Client{}
		response, error := client.Do(request)
		if error != nil {
			panic(error)
		}
		defer response.Body.Close()
		var temp player
		json.NewDecoder(response.Body).Decode(&temp)

		fmt.Println("response Status:", response.Status)
		fmt.Println("response Headers:", response.Header)
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println("response Body:", body)
		//fmt.Printf("Got %s age %d Param %s\n", tempPlayer.Name, tempPlayer.Age, val)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(temp)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}

}

func main() {
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

	r.HandleFunc("/to/{app}", MakeRequest)

	//http.ListenAndServe(":80", nil)
	log.Fatal(http.ListenAndServe(":8082", r))
}
