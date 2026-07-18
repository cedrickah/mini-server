package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
    "fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Value struct {
	Value string `json:"value"`
}

func main() {

    host := os.Getenv("DATABASE_HOST")
    port := os.Getenv("DATABASE_PORT")
    user := os.Getenv("DATABASE_USER")
    password := os.Getenv("DATABASE_PASSWORD")
    dbname := os.Getenv("DATABASE_NAME")
    
    conn :=
    fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
        host,
        port,
        user,
        password,
        dbname,
    )

	var err error

	db, err = sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}


	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS values_table(
			id SERIAL PRIMARY KEY,
			value TEXT NOT NULL
		)
	`)

	if err != nil {
		log.Fatal(err)
	}


	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/read", readHandler)
    http.HandleFunc("/health", healthHandler)


	log.Println("server listening on :8080")

	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}
}


func writeHandler(w http.ResponseWriter, r *http.Request){

	var body Value

	json.NewDecoder(r.Body).Decode(&body)


	_, err := db.Exec(
		"INSERT INTO values_table(value) VALUES($1)",
		body.Value,
	)


	if err != nil {
		http.Error(w, err.Error(),500)
		return
	}


	w.Write([]byte("saved"))
}



func readHandler(w http.ResponseWriter,r *http.Request){

	row := db.QueryRow(
		"SELECT value FROM values_table ORDER BY id DESC LIMIT 1",
	)


	var value string

	err := row.Scan(&value)

	if err != nil {
		http.Error(w,err.Error(),500)
		return
	}


	json.NewEncoder(w).Encode(
		Value{Value:value},
	)
} 

func healthHandler(w http.ResponseWriter,r *http.Request){

    w.Header().Set("Content-Type", "application/json")

    w.WriteHeader(http.StatusOK)

    w.Write([]byte(`{"status":"ok"}`))
} 
