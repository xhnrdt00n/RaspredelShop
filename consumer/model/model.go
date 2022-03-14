package model

type Category struct {
	Id     int         `json:"id"`
	Name   string      `json:"name"`
	Parent interface{} `json:"parent"`
}

type Product struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Price           string `json:"price"`
	ImageUrl        string `json:"image_url"`
	ProductCategory *int   `json:"item_category,omitempty"`
}
