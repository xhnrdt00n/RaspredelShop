package db

import "go.mongodb.org/mongo-driver/bson"

type DBService interface {
	Close() error

	GetAllCategories() ([]bson.M, error)
	GetProductsById(string) ([]bson.M, error)
}
