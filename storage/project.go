package storage

import (
	"context"
	"fmt"
	"reflect"

	"bitbucket.org/toggly/toggly-server/models"
	dbStore "github.com/nodely/go-mongo-store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type mgoProject struct {
	Storage *dbStore.DbStorage
	CRUD    dbStore.CRUD
}

func (a *mgoProject) List() []*models.Project {
	results := make([]*models.Project, 0)
	cursor, err := a.CRUD.Find(bson.D{{}}, options.Find().SetSort(bson.D{{"name", 1}}))
	if err != nil {
		fmt.Println(err.Error())
		return results
	}
	for cursor.Next(context.TODO()) {
		var rec models.Project
		cursor.Decode(&rec)
		results = append(results, &rec)
	}
	return results
}

func (a *mgoProject) Get(code string) *models.Project {
	var data models.Project
	a.CRUD.FindOne(bson.M{"code": code}).Decode(&data)
	return &data
}

func (a *mgoProject) Create(data *models.Project) (*models.Project, error) {
	// check index
	if err := a.ensureIndexes(); err != nil {
		return nil, err
	}

	ins, err := a.CRUD.Insert(data)
	if err != nil {
		return nil, err
	}
	rec, err := a.CRUD.GetItem(ins[0].(primitive.ObjectID), reflect.TypeOf(new(models.Project)))
	if err != nil {
		return nil, err
	}
	return rec.(*models.Project), nil
}

func (a *mgoProject) Update(data *models.Project) (*models.Project, error) {
	// check index
	if err := a.ensureIndexes(); err != nil {
		return nil, err
	}

	err := a.CRUD.SaveItem(data.ID, data)

	return data, err
}

func (a *mgoProject) Delete(code string) {

}

func (a *mgoProject) IsExist(code string) bool {
	return a.CRUD.Count(bson.M{"code": code}) != 0
}

func (a *mgoProject) ensureIndexes() error {
	return a.CRUD.EnsureIndexesRaw(mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "owner_id", Value: bsonx.Int32(1)},
			{Key: "code", Value: bsonx.Int32(1)},
		},
		Options: options.Index().SetUnique(true),
	})
}
