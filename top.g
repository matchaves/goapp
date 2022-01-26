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
	Club string
	Age  uint8
}

func Urlparams(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)
	w.Header().Set("Content-Type", "application/json")

	userID := -1
	var err error
	if val, ok := pathParams["userID"]; ok {
		userID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	commentID := -1
	if val, ok := pathParams["commentID"]; ok {
		commentID, err = strconv.Atoi(val)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "need a number"}`))
			return
		}
	}

	query := r.URL.Query()
	location := query.Get("location")

	w.Write([]byte(fmt.Sprintf(`{"userID": %d, "commentID": %d, "location": "%s" }`, userID, commentID, location)))
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

		fmt.Printf("Este serviço foi chamado por: %s age %d club %s\n", tempPlayer.Name, tempPlayer.Age, tempPlayer.Club)
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
		tempPlayer.Age = 8
		tempPlayer.Club = "Sucesso pelo metodo GET no /"
		fmt.Printf("Got %s age %d club %s\n", tempPlayer.Name, tempPlayer.Age, tempPlayer.Club)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tempPlayer)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}

}

func MakeRequest(w http.ResponseWriter, r *http.Request) {
	//pathParams := mux.Vars(r)
	if r.Method == "GET" {
		fmt.Println("UrlApp1:", os.Getenv("FOO"))
		fmt.Println("UrlApp2:", os.Getenv("FOO"))
		//fmt.Fprintf(w, os.Getenv("BAR"))
		//val := pathParams["app"]
		var tempPlayer player
		var temp player
		tempPlayer.Name = "33001app2"
		tempPlayer.Age = 33
		tempPlayer.Club = "request do app1"
		jsonData, _ := json.Marshal(tempPlayer)
		request, _ := http.NewRequest("POST", "http://localhost:32001/post", bytes.NewBuffer(jsonData))
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		client := &http.Client{}
		response, error := client.Do(request)
		body, _ := ioutil.ReadAll(response.Body)
		if error != nil {
			panic(error)
		}

		fmt.Println("response Status:", response.Status)
		fmt.Println("response Headers:", response.Header)
		fmt.Println("response Body:", string(body))
		decoder := json.NewDecoder(response.Body)
		err := decoder.Decode(&temp)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		//json.Unmarshal([]byte(body), &temp)
		//fmt.Printf("Got %s age %d club %s\n", tempPlayer.Name, tempPlayer.Age, val)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		//w.Write([]byte(fmt.Sprintf(`{"userID": "%s" }`, string(body))))
		//json.NewEncoder(w).Encode(people1)
		//w.Write([]byte(fmt.Sprintf(`{"userID": , "commentID":, "location": "" }`)))

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

	r.HandleFunc("/user/{userID}/comment/{commentID}", Urlparams)

	r.HandleFunc("/to/{app}", MakeRequest)

	//http.ListenAndServe(":80", nil)
	log.Fatal(http.ListenAndServe(":33001", r))
}
