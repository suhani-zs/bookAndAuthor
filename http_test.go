package httpsMithali

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	//"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAll(t *testing.T) {
	testcases := []struct {
		desc   string
		input  string
		output []Book
	}{
		{"details of all books ", "/books", []Book{
			{"XYZ", 1, &Author{1, "suhani", "siddhu", "25/04/2001", "roli"}, "arihant", "06-11-1976"},
			{"abc", 2, &Author{2, "Thomas", "alwa", "26/04/2001", "wx"}, "Penguin", "08/01/1978"}}},
	}

	for _, test := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", test.input, nil)

		GetAll(w, req)
		resp := w.Result()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var allBooks []Book

		err = json.Unmarshal(data, &allBooks)
		if err != nil {
			return
		}

		assert.Equal(t, test.output, allBooks)

		err = resp.Body.Close()

	}
}

func TestGetByID(t *testing.T) {
	testcases := []struct {
		desc   string
		input  string
		output Book
	}{
		{"The details for book XYZ: ", "/books/xyz",
			Book{"ABC", 1, &Author{1, "suhani", "siddhu", "25/04/2001", ""}, "arihant", "06-11-1976"}},
	}

	for _, test := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", test.input, nil)

		GetByID(w, req)
		resp := w.Result()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		var Book Book

		err = json.Unmarshal(data, &Book)
		if err != nil {
			return
		}

		assert.Equal(t, test.output, Book)

		err = resp.Body.Close()
	}
}

func TestPostBook(t *testing.T) {
	testcases := []struct {
		desc       string
		book       Book
		statusCode int
	}{
		{"Details posted.", Book{"ABC", 1, &Author{1, "suhani", "siddhu", "25/04/2001", "roli"}, "Arihant", "20/09/2020"},
			http.StatusCreated},
		{"Invalid book name.", Book{"", 00, &Author{4, "suhani", "siddhu", "25/04/2001", "roli"}, "Oxford", "21/04/1995"},
			http.StatusBadRequest},
	}

	for _, test := range testcases {
		newData, _ := json.Marshal(test.book)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/books", bytes.NewBuffer(newData))
		PostBook(w, req)
		resp := w.Result()

		assert.Equal(t, test.statusCode, resp.StatusCode)

	}
}

func TestPostAuthor(t *testing.T) {
	testcases := []struct {
		desc       string
		author     Author
		statusCode int
	}{
		{"Details posted.", Author{1, "suhani", "siddhu", "25/04/2001", "roli"},
			http.StatusCreated},
		{"Invalid book name.", Author{4, "suhani", "siddhu", "25/04/2001", "roli"},
			http.StatusBadRequest},
	}

	for _, test := range testcases {
		newData, _ := json.Marshal(test.author)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/author", bytes.NewBuffer(newData))
		PostBook(w, req)
		resp := w.Result()

		assert.Equal(t, test.statusCode, resp.StatusCode)

	}

}
func TestPutBook(t *testing.T) {
	testcases := []struct {
		desc      string
		reqId     string
		reqBody   Book
		expRes    Book
		expStatus int
	}{
		{"Valid Details", "1", Book{title: "hello", author1: "", publication: "Pengiun", publicationdate: "11/03/2002"}, Book{}, 200},
		{"Publication should be Scholastic/Pengiun/Arihanth", "1", Book{title: "harry potter", author1: nil, publication: "adda", publicationdate: "11/03/2002"}, Book{}, http.StatusBadRequest},
		{"Published date should be between 1880 and 2022", "1", Book{title: "", author1: nil, publication: "", publicationdate: "1/1/1870"}, Book{}, http.StatusBadRequest},
		{"Published date should be between 1880 and 2022", "1", Book{title: "", author1: nil, publication: "", publicationdate: "1/1/2222"}, Book{}, http.StatusBadRequest},
		{"Author should exist", "1", Book{}, Book{}, http.StatusBadRequest},
		{"Title can't be empty", "1", Book{title: "", author1: nil, publication: "", publicationdate: ""}, Book{}, http.StatusBadRequest},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.reqBody)
		req := httptest.NewRequest(http.MethodPost, "localhost:8000/book/"+tc.reqId, bytes.NewReader(body))
		insertBook(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
		res, _ := io.ReadAll(w.Result().Body)
		resBook := Book{}
		json.Unmarshal(res, &resBook)
		if resBook != tc.expRes {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}

func TestPutAuthor(t *testing.T) {
	testcases := []struct {
		desc      string
		reqBody   Author
		expRes    Author
		expStatus int
	}{
		{"Valid details", Author{firstName: "abc", lastName: "singh", dob: "25/04/1997", penName: "yee"}, Author{1, "abc", "singh", "25/04/1997", "yee"}, http.StatusOK},
		{"InValid details", Author{firstName: "", lastName: "", dob: "2/11/1989", penName: "Sharma"}, Author{}, http.StatusBadRequest},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.reqBody)
		req := httptest.NewRequest(http.MethodPost, "localhost:8000/author", bytes.NewReader(body))
		insertAuthor(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
		res, _ := io.ReadAll(w.Result().Body)
		resAuthor := Author{}
		json.Unmarshal(res, &resAuthor)
		if resAuthor != tc.expRes {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}
func TestDeleteBook(t *testing.T) {
	testcases := []struct {
		desc           string
		requestId      string
		expectedStatus int
	}{
		{"Valid Details", "1", http.StatusOK},
		{"Book does not exists", "90", http.StatusNotFound},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "localhost:8000/book/"+tc.requestId, nil)
		deleteBook(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expectedStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}

func TestDeleteAuthor(t *testing.T) {
	testcases := []struct {
		desc           string
		requestId      string
		expectedStatus int
	}{
		{"Valid Details", "1", http.StatusOK},
		{"Author does not exists", "90", http.StatusBadRequest},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "localhost:8000/author/"+tc.requestId, nil)
		deleteAuthor(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expectedStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}
