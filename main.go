package main

import (
	"encoding/json"
	"log"
	"net/http"

	"database/sql"
	"fmt"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Book struct (Model)
type Pharmacy struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Stock int    `json:"stock"`
	Sales string `json:"sales"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "pharmacy"
)

func connect() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Pinged")
	return db
}

// Get all drugs
func getDrugs(w http.ResponseWriter, r *http.Request) {
	var pharmacy []*Pharmacy
	db := connect()
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT * FROM drugs")
	if err != nil {
		// handle this error better than this
		log.Println(err)
		return
	}
	defer rows.Close()
	for rows.Next() {

		p := new(Pharmacy)
		switch err = rows.Scan(&p.ID, &p.Name, &p.Stock, &p.Sales); err {
		case nil:
			pharmacy = append(pharmacy, p)
		default:
			log.Println(err)
		}
	}
	json.NewEncoder(w).Encode(pharmacy)
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

}

// Get single drug
func getDrug(w http.ResponseWriter, r *http.Request) {
	db := connect()
	w.Header().Set("Content-Type", "application/json")

	sqlStatement := `SELECT * FROM drugs WHERE id=$1;`
	var pharmacy []*Pharmacy
	row := db.QueryRow(sqlStatement, 2)
	p := new(Pharmacy)
	switch err := row.Scan(&p.ID, &p.Name, &p.Stock, &p.Sales); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		pharmacy = append(pharmacy, p)
		json.NewEncoder(w).Encode(pharmacy)
	default:
		log.Println(err)
	}
}

// Add new drug
func createDrug(w http.ResponseWriter, r *http.Request) {
	db := connect()
	sqlStatement := `
INSERT INTO drugs (id,name,stock,sales)
VALUES ($1, $2, $3, $4)
RETURNING id`
	id := 0
	err := db.QueryRow(sqlStatement, 4, "Asprin", 200, "Jonathan").Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)

}

// update drug
func updateDrug(w http.ResponseWriter, r *http.Request) {
	db := connect()
	sqlStatement := `
UPDATE drugs
SET name = $2
WHERE id = $1;`
	_, err := db.Exec(sqlStatement, 1, "NewFirst")
	if err != nil {
		panic(err)
	}

	fmt.Println("Updated")

}

// delete drug
func deleteDrug(w http.ResponseWriter, r *http.Request) {
	db := connect()
	sqlStatement := `
DELETE FROM drugs
WHERE id = $1;`
	_, err := db.Exec(sqlStatement, 4)
	if err != nil {
		panic(err)
	}

	fmt.Println("deleted")

}

// Main function
func main() {
	// Init router
	r := mux.NewRouter()

	// Route handles & endpoints
	r.HandleFunc("/drugs", getDrugs).Methods("GET")
	r.HandleFunc("/drug/{id}", getDrug).Methods("GET")
	r.HandleFunc("/drugs", createDrug).Methods("POST")
	r.HandleFunc("/drugs/{id}", updateDrug).Methods("PUT")
	r.HandleFunc("/drugs/{id}", deleteDrug).Methods("DELETE")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(headers, methods, origins)(r)))

}
