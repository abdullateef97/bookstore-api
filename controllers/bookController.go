package controllers

import(
	"net/http"
	"log"
	"database/sql"
	"encoding/json"
	"strconv"
	"go-api/models"
	"go-api/helpers"
	"github.com/gorilla/mux"

)

var books []models.Book

//BookController Struct
type BookController struct{}

type addBookResponse struct {
	ID int `json:"id"`
}

//GetBooks gets All books from the DB
func (c BookController) GetBooks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		var Book models.Book
		books = []models.Book{}
	
		rows, err := db.Query("Select * from books")
		helpers.LogFatal(err)
	
		defer rows.Close()
	
		for rows.Next() {
			err := rows.Scan(&Book.ID, &Book.Author, &Book.Year, &Book.Title)
			helpers.LogFatal(err)
	
			books = append(books, Book)
		}
		json.NewEncoder(w).Encode(books)
		
	}
}

// GetBook gets details for a particular book
func (c BookController) GetBook(db *sql.DB) http.HandlerFunc {
	
	return func(w http.ResponseWriter, r *http.Request){
		var params = mux.Vars(r)
		log.Println(params)
		var id, _ = strconv.Atoi(params["id"])

		var Book models.Book
	
		rows := db.QueryRow("Select * from books where id=$1", id)
	
		err := rows.Scan(&Book.ID, &Book.Year, &Book.Title, &Book.Author)
		helpers.LogFatal(err)
	
		json.NewEncoder(w).Encode(Book)
	
	}
}

//AddBook adds new Book
func (c BookController) AddBook(db *sql.DB) http.HandlerFunc  {
	
	return func(w http.ResponseWriter, r *http.Request){
		var newBook models.Book
		var bookID int
		var addResp addBookResponse
	
		json.NewDecoder(r.Body).Decode(&newBook)
		err := db.QueryRow("insert into books (title, author, year) values($1, $2, $3) Returning id",
		 newBook.Title, newBook.Author, newBook.Year).Scan(&bookID)
		helpers.LogFatal(err)
		addResp.ID = bookID
		json.NewEncoder(w).Encode(addResp)
	}
}

//UpdateBook takes a book details and updates it
func (c BookController) UpdateBook(db *sql.DB) http.HandlerFunc  {

	return func (w http.ResponseWriter, r *http.Request){
		params := mux.Vars(r)
	
		var bookToUpdate models.Book
		json.NewDecoder(r.Body).Decode(&bookToUpdate)
	
		result, err := db.Exec("Update books set title=$1, author=$2, year=$3 where id=$4 returning id", bookToUpdate.Title, bookToUpdate.Author, bookToUpdate.Year, params["id"])
		helpers.LogFatal(err)
	
	
		rowsUpdated, err := result.RowsAffected()
		helpers.LogFatal(err)
		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

//DeleteBook to delete book by Id
func (c BookController) DeleteBook(db *sql.DB) http.HandlerFunc  {

	return func (w http.ResponseWriter, r *http.Request){
		params := mux.Vars(r)
	
		result, err  := db.Exec("delete from books where id=$1", params["id"])
		helpers.LogFatal(err)
		rowsDeleted, err := result.RowsAffected()
		helpers.LogFatal(err)
	
		json.NewEncoder(w).Encode(rowsDeleted)
	}
}