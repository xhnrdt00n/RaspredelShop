package db

import (
	"consumer/model"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

//Response example

type Category struct {
	Id            int           `bson:"id"`
	Name          string        `bson:"Name"`
	SubCategories []SubCategory `bson:"SubCategories"`
}

type SubCategory struct {
	Id int `bson:"id"`
}

func (m *MongoCon) AddCategory(cat model.Category) error {
	ctx, cancel := context.WithTimeout(m.mongoConnCtx, time.Second*10)
	defer cancel()
	coll := m.mongoConn.Database("ShopRaspredel").Collection("Categories")

	_, err := coll.InsertOne(ctx, Category{Id: cat.Id, Name: cat.Name})

	if err != nil {
		return err
	}

	addInArray := bson.D{{"$addToSet", bson.D{{"SubCategories", SubCategory{Id: cat.Id}}}}}
	_, err = coll.UpdateOne(ctx, bson.D{{"id", cat.Parent}}, addInArray)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoCon) ChangeCategory(cat model.Category) error {
	ctx, cancel := context.WithTimeout(m.mongoConnCtx, time.Second*10)
	defer cancel()
	coll := m.mongoConn.Database("ShopRaspredel").Collection("Categories")
	update := bson.D{{"$set", bson.D{{"Name", cat.Name}}}}

	//Находим и изменяем название категории
	err := coll.FindOneAndUpdate(ctx, bson.D{{"id", cat.Id}}, update)
	if err.Err() != nil {
		return err.Err()
	}

	//Ищем parent категорию
	var parent bson.D
	error := coll.FindOne(ctx, bson.D{{"SubCategories", bson.D{{"id", cat.Id}}}}).Decode(&parent)
	if error != nil {
		return error
	}

	//Проверяем что измененная не nil
	if cat.Parent == nil {
		res, err := coll.UpdateOne(ctx, bson.D{{"id", parent.Map()["id"]}}, bson.D{{"$pull", bson.D{{"SubCategories", bson.D{{"id", cat.Id}}}}}})
		if err != nil {
			return err
		}
		if res.ModifiedCount == 0 {
			return errors.New("ошибка изменения")
		}

		return nil
	}

	//Проверяем что они совпадают
	if float64(parent.Map()["id"].(int32)) == cat.Parent.(float64) {

		return nil
	}

	//Если не совпадают - удаляем
	_, error = coll.UpdateOne(ctx, bson.D{{"id", parent.Map()["id"]}}, bson.D{{"$pull", bson.D{{"SubCategories", bson.D{{"id", cat.Id}}}}}})
	if error != nil {
		return error
	}

	//Добавляем категорию новой категории отцу
	addInArray := bson.D{{"$addToSet", bson.D{{"SubCategories", SubCategory{Id: cat.Id}}}}}
	_, error = coll.UpdateOne(ctx, bson.D{{"id", cat.Parent}}, addInArray)
	if error != nil {
		return error
	}

	return nil
}

func (m *MongoCon) DeleteCategory(cat model.Category) error {
	ctx, cancel := context.WithTimeout(m.mongoConnCtx, time.Second*10)
	defer cancel()
	coll := m.mongoConn.Database("ShopRaspredel").Collection("Categories")

	//удаляем категорию
	err := coll.FindOneAndDelete(ctx, bson.D{{"id", cat.Id}})
	if err.Err() != nil {
		return err.Err()
	}

	//если у категории указан родитель, от у него удаляем категорию
	if cat.Parent != nil {
		_, error := coll.UpdateOne(ctx, bson.D{{"id", cat.Parent}}, bson.D{{"$pull", bson.D{{"SubCategories", bson.D{{"id", cat.Id}}}}}})
		if error != nil {
			return error
		}
	}

	//удаляем все айтемы у которых была эта категория //TODO обнулить категорию
	collProducts := m.mongoConn.Database("ShopRaspredel").Collection("Products")
	result, error := collProducts.DeleteMany(ctx, bson.D{{"Category", cat.Id}})
	if error != nil && result.DeletedCount > 0 {
		if err == nil {
			return errors.New("not data to delete")
		}
		return error
	}

	return nil
}
