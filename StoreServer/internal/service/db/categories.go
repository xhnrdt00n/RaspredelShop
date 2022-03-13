package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

//Response example
type Response struct {
	Categories []struct {
		Id       int    `json:"id"`
		Name     string `json:"Name"`
		Products []struct {
			Id       int    `json:"id"`
			Name     string `json:"Name"`
			ImageUrl string `json:"ImageUrl"`
			Price    string `json:"Price"`
		} `json:"Products"`
		SubCategories []struct {
			Name string `json:"Name"`
		} `json:"SubCategories"`
	} `json:"Categories"`
}

func (m *MongoCon) GetAllCategories() ([]bson.M, error) {
	ctx, cancel := context.WithTimeout(m.mongoConnCtx, time.Second*10)
	defer cancel()
	coll := m.mongoConn.Database("ShopRaspredel").Collection("Shoper")
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var responseMongo []bson.M
	err = cursor.All(ctx, &responseMongo)
	if err != nil {
		return nil, err
	}

	return responseMongo, nil
}
