package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Pallinder/go-randomdata"

	"github.com/gorilla/mux"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func main() {
	file, err := os.OpenFile("/var/log/goapp.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		ErrorLogger.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	InfoLogger.Println("Iniciando aplicacao APP2")
	r := mux.NewRouter()

	r.HandleFunc("/", PrintEnv)
	r.HandleFunc("/city", PrintEnv)

	//http.ListenAndServe(":80", nil)
	InfoLogger.Println(http.ListenAndServe(":8080", r))
}

func PrintEnv(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println("FOO:", os.Getenv("URLAPP1"))
		//fmt.Fprintf(w, os.Getenv("BAR"))

		step := fmt.Sprintf("City %v", randomdata.City())
		WarningLogger.Printf("City %v", randomdata.City())
		InfoLogger.Printf("City %v", randomdata.City())
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(step)

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed."))
	}

}
