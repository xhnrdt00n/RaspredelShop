package db

import (
	"consumer/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type Product struct {
	Id       int    `bson:"id"`
	Name     string `bson:"Name"`
	ImageUrl string `bson:"ImageUrl"`
	Price    string `bson:"Price"`
	Category *int   `bson:"Category"`
}

func (m *MongoCon) AddProduct(prod model.Product) error {
	ctx, cancel := context.WithTimeout(m.mongoConnCtx, time.Second*10)
	defer cancel()
	coll := m.mongoConn.Database("ShopRaspredel").Collection("Products")

	_, err := coll.InsertOne(ctx, Product{Id: prod.Id, Name: prod.Name, ImageUrl: prod.ImageUrl, Price: prod.Price, Category: prod.ProductCategory})

	if err != nil {
		return err
	}

	return nil
}
func (m *MongoCon) ChangeProduct(prod model.Product) error {
	ctx, cancel := context.WithTimeout(m.mongoConnCtx, time.Second*10)
	defer cancel()
	coll := m.mongoConn.Database("ShopRaspredel").Collection("Products")

	update := bson.D{{"$set", bson.D{{"Name", prod.Name}, {"ImageUrl", prod.ImageUrl}, {"Price", prod.Price}, {"Category", prod.ProductCategory}}}}

	//Находим и изменяем название категории
	err := coll.FindOneAndUpdate(ctx, bson.D{{"id", prod.Id}}, update)
	if err.Err() != nil {
		return err.Err()
	}

	return nil
}
func (m *MongoCon) DeleteProduct(prod model.Product) error {
	ctx, cancel := context.WithTimeout(m.mongoConnCtx, time.Second*10)
	defer cancel()
	coll := m.mongoConn.Database("ShopRaspredel").Collection("Products")

	err := coll.FindOneAndDelete(ctx, bson.D{{"id", prod.Id}})
	if err.Err() != nil {
		return err.Err()
	}
	return nil
}
