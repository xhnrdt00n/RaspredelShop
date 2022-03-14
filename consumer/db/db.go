package db

import (
	"consumer/model"
)

type DBService interface {
	Close() error

	AddCategory(model.Category) error
	ChangeCategory(model.Category) error
	DeleteCategory(model.Category) error
	//

	AddProduct(model.Product) error
	ChangeProduct(model.Product) error
	DeleteProduct(model.Product) error
}
