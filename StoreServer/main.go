package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/categories", handler) // each request calls handler
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
		fmt.Fprintf(w, "{\n  \"Categories\": [\n    {\n      \"id\": 1,\n      \"Name\" : \"Ноутбук\",\n      \"Products\" : [\n        {\n          \"id\":1,\n          \"Name\":\"12312\",\n          \"ImageUrl\": \"123.com\",\n          \"Price\": 125.6\n        },\n        {\n          \"id\":2,\n          \"Name\":\"Acer\",\n          \"ImageUrl\": \"https://mobile-review.com/articles/2020/image/acer-aspire-5/pic/2.jpg\",\n          \"Price\": 123.5\n        }\n      ],\n      \n      \"SubCategories\":[\n        {\n          \"Name\":\"Русские\"\n        },\n        {\n          \"Name\":\"Китайские\"\n        }\n      ]\n    },\n        {\n      \"id\" : 2,\n      \"Name\" : \"Холодильники\",\n      \"Products\" : [\n        {\n          \"id\":3,\n          \"Name\":\"New Eva\",\n          \"ImageUrl\": \"https://static.onlinetrade.ru/img/fullreviews/53166/1_big.jpg\",\n          \"Price\": 12005.6\n        },\n        {\n          \"id\":4,\n          \"Name\":\"Old Mod\",\n          \"ImageUrl\": \"https://img.mvideo.ru/Pdb/20070442b.jpg\",\n          \"Price\": 12312.5\n        }\n      ],\n      \n      \"SubCategories\":[\n        {\n          \"Name\":\"С морозилкой\"\n        },\n        {\n          \"Name\":\"Умные\"\n        }\n      ]\n    }\n  ]\n}")
	}
}
