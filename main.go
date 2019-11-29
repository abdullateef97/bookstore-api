package main
import(
	"net/http"
	"log"
	"database/sql"
	"encoding/json"
	"strconv"
	"os"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
)

type book struct {
	ID int `json:"id,omitempty"`
	Author string `json:"author"`
	Title string `json:"title"`
	Year string `json:"year"`
}

type addBookResponse struct {
	ID int `json:"id"`
}

var books []book

var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatal(err error){
	if(err != nil){
		log.Fatal(err)
	}
}

func main() {
	_, err := pq.ParseURL(os.Getenv("DB_URL"))
	logFatal(err)

	db, err = sql.Open("postgres", os.Getenv("DB_URL"))
	logFatal(err)

	err = db.Ping()
	logFatal(err)

	router := mux.NewRouter()

	router.HandleFunc("/books", getAllBooks).Methods("GET")
	router.HandleFunc("/book/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/book/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/book/{id}", deleteBook).Methods("DELETE")


	log.Fatal(http.ListenAndServe(":8000", router))
}

func getAllBooks(w http.ResponseWriter, r *http.Request){
	var Book book
	books = []book{}

	rows, err := db.Query("Select * from books")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&Book.ID, &Book.Author, &Book.Year, &Book.Title)
		logFatal(err)

		books = append(books, Book)
	}
	json.NewEncoder(w).Encode(books)
	
}

func getBook(w http.ResponseWriter, r *http.Request){
	var params = mux.Vars(r)
	log.Println(params)
	var id, _ = strconv.Atoi(params["id"])
	// var query = r.URL.Query()
	// log.Println(query["boy"][0])
	var Book book

	rows := db.QueryRow("Select * from books where id=$1", id)

	err := rows.Scan(&Book.ID, &Book.Year, &Book.Title, &Book.Author)
	logFatal(err)

	json.NewEncoder(w).Encode(Book)

}

func addBook(w http.ResponseWriter, r *http.Request){
	var newBook book
	var bookID int
	var addResp addBookResponse

	json.NewDecoder(r.Body).Decode(&newBook)
	err := db.QueryRow("insert into books (title, author, year) values($1, $2, $3) Returning id",
	 newBook.Title, newBook.Author, newBook.Year).Scan(&bookID)
	logFatal(err)
	addResp.ID = bookID
	json.NewEncoder(w).Encode(addResp)
}
func updateBook(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)

	var bookToUpdate book
	json.NewDecoder(r.Body).Decode(&bookToUpdate)

	result, err := db.Exec("Update books set title=$1, author=$2, year=$3 where id=$4 returning id", bookToUpdate.Title, bookToUpdate.Author, bookToUpdate.Year, params["id"])
	logFatal(err)


	rowsUpdated, err := result.RowsAffected()
	logFatal(err)
	json.NewEncoder(w).Encode(rowsUpdated)
}
func deleteBook(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)

	result, err  := db.Exec("delete from books where id=$1", params["id"])
	logFatal(err)
	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)
	
	log.Println("Delete Books")
}