package db

import (
	"consumer/model"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

//Response example

type Category struct {
	Id            int           `bson:"id"`
	Name          string        `bson:"name"`
	SubCategories []SubCategory `bson:"childCategories"`
	Products      []Product     `bson:"products"`
}

type SubCategory struct {
	Id   int    `bson:"id"`
	Name string `bson:"name"`
}

type Product struct {
	Id       int    `bson:"id"`
	Name     string `bson:"name"`
	ImageUrl string `bson:"image"`
	Price    string `bson:"price"`
}

//COMPLETE
func (m *MongoCon) AddCategory(cat model.Category) error {
	ctx, cancel := context.WithTimeout(m.mongoConnCtx, time.Second*10)
	defer cancel()
	coll := m.mongoConn.Database("ShopRaspredel").Collection("Categories")

	_, err := coll.InsertOne(ctx, Category{Id: cat.Id, Name: cat.Name})

	if err != nil {
		return err
	}

	var parent Category
	err = coll.FindOne(ctx, bson.D{{"id", cat.Parent}}).Decode(&parent)
	if err != nil {
		return err
	}

	if parent.SubCategories == nil {
		addInArray := bson.D{{"$set", bson.D{{"childCategories", []SubCategory{{Id: cat.Id, Name: cat.Name}}}}}}
		error := coll.FindOneAndUpdate(ctx, bson.D{{"id", cat.Parent}}, addInArray)
		if error.Err() != nil {
			return error.Err()
		}
		return nil
	}

	addInArray := bson.D{{"$addToSet", bson.D{{"childCategories", SubCategory{Id: cat.Id, Name: cat.Name}}}}}
	error := coll.FindOneAndUpdate(ctx, bson.D{{"id", cat.Parent}}, addInArray)
	if error.Err() != nil {
		return error.Err()
	}
	return nil
}

//COMPLETE
func (m *MongoCon) ChangeCategory(cat model.Category) error {
	ctx, cancel := context.WithTimeout(m.mongoConnCtx, time.Second*10)
	defer cancel()
	coll := m.mongoConn.Database("ShopRaspredel").Collection("Categories")
	update := bson.D{{"$set", bson.D{{"name", cat.Name}}}}

	//Находим и изменяем название категории
	err := coll.FindOneAndUpdate(ctx, bson.D{{"id", cat.Id}}, update)
	if err.Err() != nil {
		return err.Err()
	}

	//Ищем parent категорию
	var parent Category
	//db.users.find({ 'emails':{ $elemMatch: {'address': 'user@gmail.com'}}})
	error := coll.FindOne(ctx, bson.D{{"childCategories", bson.D{{"$elemMatch", bson.D{{"id", cat.Id}}}}}}).Decode(&parent)
	if error != nil && error != mongo.ErrNoDocuments {
		return error
	}

	//Проверяем что измененная не nil
	if cat.Parent == nil {
		res, err := coll.UpdateOne(ctx, bson.D{{"id", parent.Id}}, bson.D{{"$pull", bson.D{{"childCategories", bson.D{{"id", cat.Id}, {"name", cat.Name}}}}}})
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
		if parent.SubCategories != nil {

			for key, value := range parent.SubCategories {
				if value.Id == cat.Id {
					parent.SubCategories[key].Name = cat.Name
				}
			}

			update := bson.D{{"$set", parent}}
			n, _ := coll.UpdateOne(ctx, bson.D{{"id", parent.Id}}, update)
			if n != nil {
				fmt.Println("ok")
			}
		}

		return nil
	}

	//Если не совпадают - удаляем
	_, error = coll.UpdateOne(ctx, bson.D{{"id", parent.Id}}, bson.D{{"$pull", bson.D{{"childCategories", bson.D{{"id", cat.Id}}}}}})
	if error != nil {
		return error
	}

	//Добавляем категорию новой категории отцу

	addInArray := bson.D{{"$addToSet", bson.D{{"childCategories", SubCategory{Id: cat.Id, Name: cat.Name}}}}}
	_, error = coll.UpdateOne(ctx, bson.D{{"id", cat.Parent}}, addInArray)
	if error != nil {
		if strings.Contains(error.Error(), "Cannot apply $addToSet to non-array field") {
			addInArray := bson.D{{"$set", bson.D{{"childCategories", []SubCategory{{Id: cat.Id, Name: cat.Name}}}}}}
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

//COMPLETE
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
		_, error := coll.UpdateOne(ctx, bson.D{{"id", cat.Parent}}, bson.D{{"$pull", bson.D{{"childCategories", bson.D{{"id", cat.Id}}}}}})
		if error != nil {
			return error
		}
	}

	//удаляем все айтемы у которых была эта категория //TODO обнулить категорию
	//collProducts := m.mongoConn.Database("ShopRaspredel").Collection("Products")
	//result, error := collProducts.DeleteMany(ctx, bson.D{{"Category", cat.Id}})
	//if error != nil && result.DeletedCount > 0 {
	//	if err == nil {
	//		return errors.New("not data to delete")
	//	}
	//	return error
	//}

	return nil
}

//COMPLETE
func (m *MongoCon) AddProduct(prod model.Product) error {
	ctx, cancel := context.WithTimeout(m.mongoConnCtx, time.Second*10)
	defer cancel()
	coll := m.mongoConn.Database("ShopRaspredel").Collection("Categories")

	if prod.ProductCategory == nil {
		return errors.New("no category in product")
	}

	var category Category
	err := coll.FindOne(ctx, bson.D{{"id", prod.ProductCategory}}).Decode(&category)

	category.Products = append(category.Products, Product{Id: prod.Id, Name: prod.Name, ImageUrl: prod.ImageUrl, Price: prod.Price})

	update := bson.D{{"$set", category}}
	n, _ := coll.UpdateOne(ctx, bson.D{{"id", prod.ProductCategory}}, update)
	if n != nil {
		fmt.Println("ok")
	}

	if err != nil {
		return err
	}

	return nil
}

//COMPLETE
func (m *MongoCon) ChangeProduct(prod model.Product) error {
	ctx, cancel := context.WithTimeout(m.mongoConnCtx, time.Minute*10)
	defer cancel()
	coll := m.mongoConn.Database("ShopRaspredel").Collection("Categories")

	//find parent
	var productOwner Category
	//db.users.find({ 'emails':{ $elemMatch: {'address': 'user@gmail.com'}}})
	error := coll.FindOne(ctx, bson.D{{"products", bson.D{{"$elemMatch", bson.D{{"id", prod.Id}}}}}}).Decode(&productOwner)
	if error != nil && error != mongo.ErrNoDocuments {
		return error
	}

	//assert parent with category changed
	if productOwner.Id == *prod.ProductCategory {
		for key, value := range productOwner.Products {
			if value.Id == prod.Id {
				productOwner.Products[key].Name = prod.Name
				productOwner.Products[key].Price = prod.Price
				productOwner.Products[key].ImageUrl = prod.ImageUrl
			}
		}

		//change parent and out
		update := bson.D{{"$set", productOwner}}
		n, _ := coll.UpdateOne(ctx, bson.D{{"id", productOwner.Id}}, update)
		if n != nil {
			fmt.Println("ok")
		}
		return nil
	}

	//delete from old category
	var keyProduct int
	for key, value := range productOwner.Products {
		if value.Id == prod.Id {
			keyProduct = key
		}
	}
	productOwner.Products = append(productOwner.Products[:keyProduct], productOwner.Products[keyProduct+1:]...)
	updateOldCategory := bson.D{{"$set", productOwner}}
	res, _ := coll.UpdateOne(ctx, bson.D{{"id", productOwner.Id}}, updateOldCategory)
	if res.ModifiedCount > 1 {
		fmt.Println("delete from parent")
	}

	//Add to new category
	var newCategoryOwner Category
	errNewCat := coll.FindOne(ctx, bson.D{{"id", prod.ProductCategory}}).Decode(&newCategoryOwner)

	newCategoryOwner.Products = append(newCategoryOwner.Products, Product{Id: prod.Id, Name: prod.Name, ImageUrl: prod.ImageUrl, Price: prod.Price})

	updateNewCat := bson.D{{"$set", newCategoryOwner}}
	n, _ := coll.UpdateOne(ctx, bson.D{{"id", newCategoryOwner.Id}}, updateNewCat)
	if n != nil {
		fmt.Println("ok")
	}

	if errNewCat != nil {
		return errNewCat
	}

	return nil
}

//COMPLETE
func (m *MongoCon) DeleteProduct(prod model.Product) error {
	ctx, cancel := context.WithTimeout(m.mongoConnCtx, time.Second*10)
	defer cancel()
	coll := m.mongoConn.Database("ShopRaspredel").Collection("Categories")

	var category Category
	err := coll.FindOne(ctx, bson.D{{"id", prod.ProductCategory}}).Decode(&category)
	if err != nil {
		return err
	}
	var keyProduct int
	for key, value := range category.Products {
		if value.Id == prod.Id {
			keyProduct = key
		}
	}
	category.Products = append(category.Products[:keyProduct], category.Products[keyProduct+1:]...)
	update := bson.D{{"$set", category}}
	res, _ := coll.UpdateOne(ctx, bson.D{{"id", prod.ProductCategory}}, update)
	if res.ModifiedCount > 1 {
		fmt.Println("ok")
	}
	return nil
}
