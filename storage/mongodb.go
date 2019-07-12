package storage

import (
	dbStore "github.com/nodely/go-mongo-store"
)

// Storage struct
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
