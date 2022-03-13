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
		fmt.Fprintf(w, "[{\"Name\":\"privetpoka\",\"Products\":null,\"SubCategories\":null,\"_id\":\"622e676c066b4bf39fc2047b\",\"id\":40},{\"Name\":\"Ноутбук\",\"Products\":[{\"ImageUrl\":\"123.com\",\"Name\":\"12312\",\"Price\":125.6,\"id\":1},{\"ImageUrl\":\"https://mobile-review.com/articles/2020/image/acer-aspire-5/pic/2.jpg\",\"Name\":\"Acer\",\"Price\":123.5,\"id\":2}],\"SubCategories\":[{\"Name\":\"Русские\"},{\"Name\":\"Китайские\"}],\"_id\":\"622e692ee0f80bdd6793d527\",\"id\":1},{\"Name\":\"Холодильники\",\"Products\":[{\"ImageUrl\":\"https://static.onlinetrade.ru/img/fullreviews/53166/1_big.jpg\",\"Name\":\"New Eva\",\"Price\":12005.6,\"id\":3},{\"ImageUrl\":\"https://img.mvideo.ru/Pdb/20070442b.jpg\",\"Name\":\"Old Mod\",\"Price\":12312.5,\"id\":4}],\"SubCategories\":[{\"Name\":\"С морозилкой\"},{\"Name\":\"Умные\"}],\"_id\":\"622e6970e0f80bdd679453c5\",\"id\":2}]")
	}
}
