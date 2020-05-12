package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	// Import GORM-related packages.
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"

    "gitlab.iqt.org/rashley/covid-test-db/models/poc"
)

func getDB () *gorm.DB {
	const addr = "postgresql://covid_bug@roach:26257/covid_diagnostics?sslmode=disable"
    db, err := gorm.Open("postgres", addr)
    if err != nil {
        log.Fatal(err)
    }
    //defer db.Close()

    // Set to `true` and GORM will print out all DB queries.
    db.LogMode(false)

    return db
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	log.Println("homelink called")
	fmt.Fprintf(w, "Welcome home!")
}

func createPoc(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	defer db.Close()
	var p poc.Poc

	err := json.NewDecoder(r.Body).Decode(&p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	created, err := poc.Create(db, p.Name, p.Email, p.Phone)
	if err != nil {
        log.Print(err)
    }
	json.NewEncoder(w).Encode(created)
}

func updatePoc(w http.ResponseWriter, r *http.Request) {
	pocID, err := uuid.Parse(mux.Vars(r)["id"])
	db := getDB()
	defer db.Close()
	var p poc.Poc

	err = json.NewDecoder(r.Body).Decode(&p)

    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    p.Id = pocID

	created, err := poc.Update(db, &p)
	if err != nil {
        log.Print(err)
    }
	json.NewEncoder(w).Encode(created)
}

func getPocList(w http.ResponseWriter, r *http.Request) {
	db := getDB()
	defer db.Close()

	results, err := poc.FetchList(db)
	if err != nil {
        log.Print(err)
    }
	json.NewEncoder(w).Encode(results)
}

func getPoc(w http.ResponseWriter, r *http.Request) {
	pocID, err := uuid.Parse(mux.Vars(r)["id"])
	db := getDB()
	defer db.Close()

	result, err := poc.FetchById(db, pocID)
	if err != nil {
        log.Print(err)
    }
	json.NewEncoder(w).Encode(result)
}

func main() {
	logFile, err := os.OpenFile("api.log", os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
	if err != nil {
	    panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.Println("Logging Started")

	router := mux.NewRouter().StrictSlash(false)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router.HandleFunc("/", homeLink)
	router.HandleFunc("/pocs", getPocList).Methods("GET", "OPTIONS")
	router.HandleFunc("/pocs", createPoc).Methods("POST")
	router.HandleFunc("/pocs/{id}", getPoc).Methods("GET")
	router.HandleFunc("/pocs/{id}", updatePoc).Methods("PUT")
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(headersOk, originsOk, methodsOk)(router)))
}