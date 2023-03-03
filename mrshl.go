package main

import (
	"encoding/json"
	"fmt"
)

// /////////////////////////////////////////////////////////////////////////////////////////////
type EmailData struct {
	Email   string `json:"email"`
	Content string `json:"content"`
	Text    string `json:"text"`
}

// //////////////////////////////////////////////////////////////////////////////////////////////
type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	City  string `json:"city"`
	Books []Book `json:"books"`
}

type Book struct {
	Name string `json:"name"`
	Year int    `json:"year"`
}

func main() {
	byt := []byte(`{"name":"John","age":20,"city":"London","books":[{"name":"Bookname","year":1990},{"name":"Bookname2","year":2090}]}`)
	msg := []byte()
	var dat User

	if err := json.Unmarshal(byt, &dat); err != nil {
		panic(err)
	}
	fmt.Println(dat)
	//fmt.Println(dat["name"])
	//fmt.Println(dat["books"].([]interface{})[0].(map[string]interface{})["name"])
}

func Serialize() {
	var books []Book
	book1 := Book{
		Name: "Bookname",
		Year: 1990,
	}
	book2 := Book{
		Name: "Bookname2",
		Year: 2090,
	}
	books = append(books, book1, book2)
	sv := User{
		Name:  "John",
		Age:   20,
		City:  "London",
		Books: books,
	} //map[string]interface{}{"field1": "value", "field2": 123, "field3": true, "arr": []string{"a", "b", "c"}}
	BoolVar, _ := json.Marshal(sv)
	fmt.Println(string(BoolVar))
}
