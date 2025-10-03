package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-postgres/models"
	"log"
	"net/http"
	"os"
	"strconv"
	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type response struct {
	ID int64 `json:"id,omitempty"`
	// Status string `json:"status"`
	Message string `json:"message,omitempty"`
}

func CreateConnection() *sql.DB {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	fmt.Println("before panic1")
	if err != nil {
		panic(err)
	}
	fmt.Println("before ping")
	err = db.Ping()

	fmt.Println("before panic2")
	if err != nil {
		panic(err)
	}
	fmt.Println("after panic 2")
	fmt.Println("Successfully connected to database...")

	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {

	var stock models.Stock
	// w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&stock)

	fmt.Println(stock)

	if err != nil {
		log.Fatalf("Unable to Decode request body %v", err)
	}

	insertID := insertStock(stock)

	res := response{
		ID: insertID,
		// Status: "success",
		Message: "Stock Created Successfully",
	}

	json.NewEncoder(w).Encode(res)

}

func Getstock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string to int %v", err)
	}

	stock, err := getStock(int64(id))

	if err != nil {
		log.Fatalf("Unable to get stock %v", err)
	}

	json.NewEncoder(w).Encode(stock)

}

func GetAllstocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := getAllStocks()

	if err != nil {
		log.Fatalf("Unable to get all stocks %v", err)
	}

	json.NewEncoder(w).Encode(stocks)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string to int %v", err)
	}

	var stock models.Stock

	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Unable to Decode request body %v", err)
	}

	updatedRow := updateStock(int64(id), stock)

	msg := fmt.Sprintf("Stock updated successfully. Total rows/record affected %v", updatedRow)

	res := response{
		ID: int64(id),
		// Status: "success",
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)

}

func DeleteStock(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string to int %v", err)
	}

	deletedRows := deleteStock(int64(id))

	msg := fmt.Sprintf("Stock deleted successfully. Total rows/record affected %v", deletedRows)

	res := response{
		ID: int64(id),
		// Status: "success",
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}


func insertStock(stock models.Stock) int64 {
	db:= CreateConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO stocks (name,price,company) VALUES ($1,$2,$3) RETURNING stockid`
	var id int64

	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)
	return id
}

func getStock(id int64) (models.Stock, error) {
	db:= CreateConnection()
	defer db.Close()

	var stock models.Stock

	sqlStatement := `SELECT * FROM stocks WHERE stock_id=$1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("Unable to scan the row %v", err)
	}


	return stock, err




}

func getAllStocks() ([]models.Stock, error) {
	db:= CreateConnection()
	defer db.Close()

	var stocks []models.Stock

	sqlStatement := `SELECT * FROM stocks`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var stock models.Stock

		err := rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

		if err != nil {
			log.Fatalf("Unable to scan the row %v", err)
		}

		stocks = append(stocks, stock)
	}

	return stocks, err
}


func updateStock(id int64, stock models.Stock) int64 {
	db:= CreateConnection()
	defer db.Close()	

	sqlStatement := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stock_id=$1`

	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)

	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error occured while checking the affected rows %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)
	return rowsAffected
}

func deleteStock(id int64) int64 {
	db:= CreateConnection()
	defer db.Close()

	sqlStatement := `DELETE FROM stocks WHERE stock_id=$1`

	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}

	affect, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error occured while checking the affected rows %v", err)
	}

	fmt.Printf("Total rows/record affected %v", affect)
	return affect
}


