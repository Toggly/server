package storage

import (
	dbStore "github.com/nodely/go-mongo-store"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoStorage struct
type MongoStorage struct {
	Dbs *dbStore.DbStorage
}

// GetProjectsCollection func
func (db *MongoStorage) GetProjectsCollection() dbStore.CRUD {
	return db.Dbs.GetDbCollection("projects")
}

// GetEnvsCollection func
func (db *MongoStorage) GetEnvsCollection() dbStore.CRUD {
	return db.Dbs.GetDbCollection("envs")
}

// GetEnvsKeysCollection func
func (db *MongoStorage) GetEnvsKeysCollection() dbStore.CRUD {
	return db.Dbs.GetDbCollection("envs_keys")
}

// GetPackagesCollection func
func (db *MongoStorage) GetPackagesCollection() dbStore.CRUD {
	return db.Dbs.GetDbCollection("packages")
}

// GetObjectsCollection func
func (db *MongoStorage) GetObjectsCollection() dbStore.CRUD {
	return db.Dbs.GetDbCollection("objects")
}

// GetParamsCollection func
func (db *MongoStorage) GetParamsCollection() dbStore.CRUD {
	return db.Dbs.GetDbCollection("params")
}

// ProjectCRUD func
func (db *MongoStorage) ProjectCRUD() Project {
	return &mgoProject{Storage: db.Dbs, CRUD: db.GetProjectsCollection()}
}

// ParamCRUD func
func (db *MongoStorage) ParamCRUD(project primitive.ObjectID) Parameter {
	return &mgoParams{Storage: db.Dbs, CRUD: db.GetParamsCollection(), ProjectID: project}
}

// EnvCRUD func
func (db *MongoStorage) EnvCRUD(project primitive.ObjectID) Environment {
	return &mgoEnvs{Storage: db.Dbs, CRUD: db.GetEnvsCollection(), ProjectID: project}
}

// EnvKeyCRUD func
func (db *MongoStorage) EnvKeyCRUD(project primitive.ObjectID) EnvironmentKey {
	return &mgoEnvsKeys{Storage: db.Dbs, CRUD: db.GetEnvsKeysCollection(), ProjectID: project}
}

// PackageCRUD func
func (db *MongoStorage) PackageCRUD(project primitive.ObjectID) Package {
	return &mgoPackage{Storage: db.Dbs, CRUD: db.GetPackagesCollection(), ProjectID: project}
}
