package db

import (
	"consumer/model"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

//Response example

type Category struct {
	Id            int           `bson:"id"`
	Name          string        `bson:"Name"`
	SubCategories []SubCategory `bson:"SubCategories"`
}

type SubCategory struct {
	Id   int    `bson:"id"`
	Name string `bson:"name"`
}

type SubCategors struct {
	Categ []SubCategory `bson:"SubCategories"`
}

func (m *MongoCon) AddCategory(cat model.Category) error {
	ctx, cancel := context.WithTimeout(m.mongoConnCtx, time.Second*10)
	defer cancel()
	coll := m.mongoConn.Database("ShopRaspredel").Collection("Categories")

	_, err := coll.InsertOne(ctx, Category{Id: cat.Id, Name: cat.Name})

	if err != nil {
		return err
	}

	var categ Category
	err = coll.FindOne(ctx, bson.D{{"id", cat.Parent}}).Decode(&categ)
	if err != nil {
		return err
	}

	if categ.SubCategories == nil {
		addInArray := bson.D{{"$set", bson.D{{"SubCategories", []SubCategory{{Id: cat.Id, Name: cat.Name}}}}}}
		error := coll.FindOneAndUpdate(ctx, bson.D{{"id", cat.Parent}}, addInArray)
		if error.Err() != nil {
			return error.Err()
		}
		return nil
	}

	addInArray := bson.D{{"$addToSet", bson.D{{"SubCategories", SubCategory{Id: cat.Id, Name: cat.Name}}}}}
	error := coll.FindOneAndUpdate(ctx, bson.D{{"id", cat.Parent}}, addInArray)
	if error.Err() != nil {
		return error.Err()
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
	var parent Category
	//db.users.find({ 'emails':{ $elemMatch: {'address': 'user@gmail.com'}}})
	error := coll.FindOne(ctx, bson.D{{"SubCategories", bson.D{{"$elemMatch", bson.D{{"id", cat.Id}}}}}}).Decode(&parent)
	if error != nil && error != mongo.ErrNoDocuments {
		return error
	}

	//Проверяем что измененная не nil
	if cat.Parent == nil {
		res, err := coll.UpdateOne(ctx, bson.D{{"id", parent.Id}}, bson.D{{"$pull", bson.D{{"SubCategories", bson.D{{"id", cat.Id}, {"name", cat.Name}}}}}})
		if err != nil {
			return err
		}
		if res.ModifiedCount == 0 {
			return errors.New("ошибка изменения")
		}

		return nil
	}

	//Проверяем что они совпадают
	if float64(parent.Id) == cat.Parent.(float64) {

		error := coll.FindOneAndUpdate(ctx, bson.D{{"SubCategories", bson.D{{"$elemMatch", bson.D{{"id", cat.Id}}}}}}, bson.D{{"$set", bson.D{{"name", cat.Name}}}})
		if error.Err() != nil {
			return error.Err()
		}

		return nil
	}

	//Если не совпадают - удаляем
	_, error = coll.UpdateOne(ctx, bson.D{{"id", parent.Id}}, bson.D{{"$pull", bson.D{{"SubCategories", bson.D{{"id", cat.Id}}}}}})
	if error != nil {
		return error
	}

	//Добавляем категорию новой категории отцу

	addInArray := bson.D{{"$addToSet", bson.D{{"SubCategories", SubCategory{Id: cat.Id, Name: cat.Name}}}}}
	_, error = coll.UpdateOne(ctx, bson.D{{"id", cat.Parent}}, addInArray)
	if error != nil {
		if strings.Contains(error.Error(), "Cannot apply $addToSet to non-array field") {
			addInArray := bson.D{{"$set", bson.D{{"SubCategories", []SubCategory{{Id: cat.Id, Name: cat.Name}}}}}}
			error := coll.FindOneAndUpdate(ctx, bson.D{{"id", cat.Parent}}, addInArray)
			if error.Err() != nil {
				return error.Err()
			}
			return nil
		}
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
