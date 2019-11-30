package main
import(
	"net/http"
	"log"
	"database/sql"

	"go-api/drivers"
	"go-api/controllers"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)


var db *sql.DB

func init() {
	gotenv.Load()
}



func main() {
	db = drivers.ConnectDB()

	router := mux.NewRouter()

	bookControllers := controllers.BookController{}

	router.HandleFunc("/books", bookControllers.GetBooks(db)).Methods("GET")
	router.HandleFunc("/book/{id}", bookControllers.GetBook(db)).Methods("GET")
	router.HandleFunc("/books", bookControllers.AddBook(db)).Methods("POST")
	router.HandleFunc("/book/{id}", bookControllers.UpdateBook(db)).Methods("PUT")
	router.HandleFunc("/book/{id}", bookControllers.DeleteBook(db)).Methods("DELETE")


	log.Fatal(http.ListenAndServe(":8000", router))
}
