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

type mgoEnvs struct {
	Storage   *dbStore.DbStorage
	CRUD      dbStore.CRUD
	ProjectID primitive.ObjectID
}

type mgoEnvsKeys struct {
	Storage   *dbStore.DbStorage
	CRUD      dbStore.CRUD
	ProjectID primitive.ObjectID
}

func (a *mgoEnvs) List() []*models.Environment {
	results := make([]*models.Environment, 0)
	cursor, err := a.CRUD.Find(bson.M{"project_id": a.ProjectID}, options.Find().SetSort(bson.D{{"protected", -1}, {"code", 1}}))
	if err != nil {
		fmt.Println(err.Error())
		return results
	}
	for cursor.Next(context.TODO()) {
		var rec models.Environment
		cursor.Decode(&rec)
		results = append(results, &rec)
	}
	return results
}

func (a *mgoEnvs) Get(code string) (*models.Environment, error) {
	var data models.Environment
	res := a.CRUD.FindOne(bson.M{"code": code, "project_id": a.ProjectID})
	res.Decode(&data)
	if res.Err() != nil {
		return nil, res.Err()
	}
	return &data, nil
}

func (a *mgoEnvs) Create(data *models.Environment) (*models.Environment, error) {
	// check index
	if err := a.ensureIndexes(); err != nil {
		return nil, err
	}

	ins, err := a.CRUD.Insert(data)
	if err != nil {
		return nil, err
	}
	rec, err := a.CRUD.GetItem(ins[0].(primitive.ObjectID), reflect.TypeOf(new(models.Environment)))
	if err != nil {
		return nil, err
	}
	return rec.(*models.Environment), nil
}

func (a *mgoEnvs) Update(data *models.Environment) (*models.Environment, error) {
	// check index
	if err := a.ensureIndexes(); err != nil {
		return nil, err
	}

	err := a.CRUD.SaveItem(data.ID, data)

	return data, err
}

func (a *mgoEnvs) Delete(code string) {}

func (a *mgoEnvs) IsExist(code string) bool {
	return a.CRUD.Count(bson.M{"code": code}) != 0
}

func (a *mgoEnvs) ensureIndexes() error {
	return a.CRUD.EnsureIndexesRaw(mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "project_id", Value: bsonx.Int32(1)},
			{Key: "code", Value: bsonx.Int32(1)},
		},
		Options: options.Index().SetUnique(true),
	})
}

func (a *mgoEnvsKeys) Provision(data *models.EnvAPIKey) (*models.EnvAPIKey, error) {
	a.ensureIndexes()

	ins, err := a.CRUD.Insert(data)
	if err != nil {
		return nil, err
	}
	rec, err := a.CRUD.GetItem(ins[0].(primitive.ObjectID), reflect.TypeOf(new(models.EnvAPIKey)))
	if err != nil {
		return nil, err
	}
	return rec.(*models.EnvAPIKey), nil
}

func (a *mgoEnvsKeys) Check(key string, secret string) error {
	var data models.EnvAPIKey
	res := a.CRUD.FindOne(bson.M{"key": key, "secret": secret})
	// todo decide do we need api key data
	res.Decode(&data)
	return res.Err()
}

func (a *mgoEnvsKeys) ensureIndexes() error {
	return a.CRUD.EnsureIndexesRaw(mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "key", Value: bsonx.Int32(1)},
			{Key: "env_id", Value: bsonx.Int32(1)},
		},
		Options: options.Index().SetUnique(true),
	})
}
