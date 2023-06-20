package mongohelper

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type InsertOneProps struct {
	Ctx            context.Context
	CollectionName string
	InsertData     interface{}
	Opts           options.InsertOneOptions
}

type FindProps struct {
	Ctx            context.Context
	CollectionName string
	Filter         interface{}
	Opts           options.FindOptions
}

type FindOneProps struct {
	Ctx            context.Context
	CollectionName string
	Filter         interface{}
	Opts           options.FindOneOptions
}

type UpdateOneProps struct {
	Ctx                context.Context
	CollectionName     string
	Filter, UpdateData interface{}
	Opts               options.UpdateOptions
}

type DeleteOneProps struct {
	Ctx            context.Context
	CollectionName string
	Filter         interface{}
	Opts           options.DeleteOptions
}
