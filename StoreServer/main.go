package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/categories/{id}", handler2)
	r.HandleFunc("/categories", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Response struct {
	Categories []struct {
		Products []struct {
			Name     string `json:"name"`
			ImageUrl string `json:"imageUrl"`
		} `json:"products"`
	} `json:"categories"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		resp := new(Response)
		resp.Categories = append(resp.Categories)
		fmt.Fprintf(w, "[{\n  \"id\":1,\n  \"Name\":\"Samsung\",\n  \"ImageUrl\":\"123.com\",\n  \"Price\": 125.6,\n  \"Category\": 3\n},\n{\n  \"id\":2,\n  \"Name\":\"Huawei\",\n  \"ImageUrl\":\"1235.com\",\n  \"Price\": 12523.6,\n  \"Category\": 3\n}]\n")
	}
}

func handler2(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		vars := mux.Vars(r)
		_, ok := vars["id"]
		if !ok {
			fmt.Println("id is missing in parameters")
		}
		fmt.Fprintf(w, "[\n  {\n    \"id\": 1,\n    \"Name\": \"Ноутбуки\",\n    \"SubCategories\": [\n      {\n        \"id\": 2\n      }\n    ]\n  },\n  {\n    \"id\": 2,\n    \"Name\": \"Адаптеры питания\",\n    \"SubCategories\": null\n  },\n  {\n    \"id\": 3,\n    \"Name\": \"Телефоны\",\n    \"SubCategories\": null\n  }\n]")
	}
}
