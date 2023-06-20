package mongohelper

import (
	"context"
	"sync"

	"appota/web-builder/config"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBHelper interface {
	GetClient() *mongo.Client
	GetCollection(name string) *mongo.Collection
	Connect(uri string) error
	Disconnect()
	InsertOne(props InsertOneProps) (*mongo.InsertOneResult, error)
	UpdateOne(props UpdateOneProps) (*mongo.UpdateResult, error)
	Find(props FindProps) (*mongo.Cursor, error)
	DeleteOne(props DeleteOneProps) (*mongo.DeleteResult, error)
	FindOne(props FindOneProps) *mongo.SingleResult
}

type helper struct {
	client *mongo.Client
}

var once sync.Once

func (h *helper) Connect(uri string) (err error) {
	once.Do(func() {
		opts := options.Client().ApplyURI(uri)
		h.client, err = mongo.Connect(context.TODO(), opts)
		if err != nil {
			return
		}

		log.Info("connected to the database successfully.")
	})

	return nil
}

func (h *helper) Disconnect() {
	log.Info("disconnecting the database...")
	h.client.Disconnect(context.TODO())
}

func (h *helper) InsertOne(props InsertOneProps) (*mongo.InsertOneResult, error) {
	col := h.GetCollection(props.CollectionName)
	res, err := col.InsertOne(props.Ctx, props.InsertData, &props.Opts)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return res, nil
}

func (h *helper) Find(props FindProps) (*mongo.Cursor, error) {
	col := h.GetCollection(props.CollectionName)
	res, err := col.Find(props.Ctx, props.Filter, &props.Opts)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return res, nil
}

func (h *helper) FindOne(props FindOneProps) *mongo.SingleResult {
	col := h.GetCollection(props.CollectionName)
	res := col.FindOne(props.Ctx, props.Filter, &props.Opts)

	return res
}

func (h *helper) UpdateOne(props UpdateOneProps) (*mongo.UpdateResult, error) {
	col := h.GetCollection(props.CollectionName)
	res, err := col.UpdateOne(props.Ctx, props.Filter, props.UpdateData, &props.Opts)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return res, nil
}

func (h *helper) DeleteOne(props DeleteOneProps) (*mongo.DeleteResult, error) {
	col := h.GetCollection(props.CollectionName)
	res, err := col.DeleteOne(props.Ctx, props.Filter, &props.Opts)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return res, nil
}

func (h *helper) GetClient() *mongo.Client {
	client := h.client
	return client
}

func (h *helper) GetCollection(name string) *mongo.Collection {
	col := h.client.Database(config.GetDatabaseName()).Collection(name)

	return col
}

func NewHelper() MongoDBHelper {
	h := &helper{}

	return h
}
