package main

import (
	"encoding/json"
	"fmt"
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

		fmt.Printf("Got %s age %d club %s\n", tempPlayer.Name, tempPlayer.Age, tempPlayer.Club)
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
		fmt.Fprintf(w, os.Getenv("BAR"))
		var tempPlayer player
		tempPlayer.Name = "APP1"
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
	pathParams := mux.Vars(r)
	if r.Method == "GET" {
		fmt.Println("UrlApp1:", os.Getenv("FOO"))
		fmt.Println("UrlApp2:", os.Getenv("FOO"))
		fmt.Fprintf(w, os.Getenv("BAR"))
		val := pathParams["app"]
		var tempPlayer player
		tempPlayer.Name = val
		tempPlayer.Age = 8
		tempPlayer.Club = "Sucesso pelo metodo GET no / MakeRequest"
		fmt.Printf("Got %s age %d club %s\n", tempPlayer.Name, tempPlayer.Age, val)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(tempPlayer)

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
	log.Fatal(http.ListenAndServe(":32001", r))
}
