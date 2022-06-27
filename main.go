package httpsMithali

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

type Author struct {
	authorId  int    `json:"authorId"`
	firstName string `json:"firstName"`
	lastName  string `json:"lastName"`
	dob       string `json:"dob"`
	penName   string `json:"penName"`
}

type Book struct {
	title           string  `json:"title"`
	bookId          int     `json:"bookId"`
	author1         *Author `json:"author1"`
	publication     string  `json:"publication"`
	publicationdate string  `json:"publicationdate"`
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "raramuri"
	dbPass := "Suhani@123"
	dbName := "public"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp(127.0.0.1:3306)/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func GetAll(response http.ResponseWriter, request *http.Request) {
	db := dbConn()
	defer db.Close()
	title := request.URL.Query().Get("title")
	includeAuthor := request.URL.Query().Get("includeAuthor")
	var rows *sql.Rows
	var err error
	if title == "" {
		rows, err = db.Query("select * from book;")
	} else {
		rows, err = db.Query("select * from book where title=?;", title)
	}
	if err != nil {
		log.Print(err)
	}
	books := []Book{}
	for rows.Next() {
		book := Book{}
		err = rows.Scan(&book.bookId, &book.title, &book.publication, &book.publicationdate, &book.author1)
		if err != nil {
			log.Print(err)
		}
		if includeAuthor == "true" {
			row := db.QueryRow("select * from Authors where id=?", book.author1.authorId)
			row.Scan(&book.author1.authorId, &book.author1.firstName, &book.author1.lastName, &book.author1.dob, &book.author1.penName)
		}
		books = append(books, book)
	}
	json.NewEncoder(response).Encode(books)
}
func GetById(response http.ResponseWriter, request *http.Request) {
	json.NewEncoder(response).Encode(Book{
		title: "jk",
		author1: &Author{
			authorId:  1,
			firstName: "suhani",
			lastName:  "siddhu",
			dob:       "25/04/2001",
			penName:   "roli",
		},
		publication:     "",
		publicationdate: "",
	})
}
