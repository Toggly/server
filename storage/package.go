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

type mgoPackage struct {
	Storage   *dbStore.DbStorage
	CRUD      dbStore.CRUD
	ProjectID primitive.ObjectID
}

func (a *mgoPackage) List() []*models.Package {
	results := make([]*models.Package, 0)
	cursor, err := a.CRUD.Find(bson.M{"project_id": a.ProjectID}, options.Find().SetSort(bson.D{{"code", 1}}))
	if err != nil {
		fmt.Println(err.Error())
		return results
	}
	for cursor.Next(context.TODO()) {
		var rec models.Package
		cursor.Decode(&rec)
		results = append(results, &rec)
	}
	return results
}

func (a *mgoPackage) Get(code string) (*models.Package, error) {
	var data models.Package
	res := a.CRUD.FindOne(bson.M{"code": code, "project_id": a.ProjectID})
	if res.Err() != nil {
		return nil, res.Err()
	}
	res.Decode(&data)
	return &data, nil
}

func (a *mgoPackage) Create(data *models.Package) (*models.Package, error) {
	// check index
	if err := a.ensureIndexes(); err != nil {
		return nil, err
	}

	ins, err := a.CRUD.Insert(data)
	if err != nil {
		return nil, err
	}
	rec, err := a.CRUD.GetItem(ins[0].(primitive.ObjectID), reflect.TypeOf(new(models.Package)))
	if err != nil {
		return nil, err
	}
	return rec.(*models.Package), nil
}

func (a *mgoPackage) Update(data *models.Package) (*models.Package, error) {
	// check index
	if err := a.ensureIndexes(); err != nil {
		return nil, err
	}

	err := a.CRUD.SaveItem(data.ID, data)

	return data, err
}

func (a *mgoPackage) Delete(code string) {}

func (a *mgoPackage) IsExist(code string) bool {
	return a.CRUD.Count(bson.M{"code": code}) != 0
}

func (a *mgoPackage) ensureIndexes() error {
	return a.CRUD.EnsureIndexesRaw(mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "project_id", Value: bsonx.Int32(1)},
			{Key: "code", Value: bsonx.Int32(1)},
		},
		Options: options.Index().SetUnique(true),
	})
}
