package media

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"

	"appota/web-builder/config"
	"appota/web-builder/db/mongohelper"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const metadataCollectionName = "metadata"

type Repository interface {
	SaveMedia(file []byte, filename string) error
	GetMedia(path string) ([]byte, string, error)
	InsertMediaData(data *Media) error
	FindMediaData(filter primitive.D, opts options.FindOptions) (*mongo.Cursor, error)
	UpdateMediaData(filter primitive.D, data *UpdateMedia) (*mongo.UpdateResult, error)
	DeleteMediaData(filter primitive.D) (*mongo.DeleteResult, error)
}

func (r *repo) SaveMedia(fileByte []byte, filename string) error {
	path := filepath.Join(config.GetStorePath(), "files")
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	path = filepath.Join(path, filename)
	dist, err := os.Create(path)
	if err != nil {
		log.Error(err)
		return err
	}
	defer dist.Close()

	fileSource := bytes.NewReader(fileByte)
	if _, err := io.Copy(dist, fileSource); err != nil {
		return err
	}
	return nil
}

func (r *repo) GetMedia(filename string) (file []byte, fname string, err error) {
	path := filepath.Join(config.GetStorePath(), "files")
	files, err := os.ReadDir(path)
	if err != nil {
		log.Error(err)
		return nil, "", err
	}

	wd := filepath.Join(path, filename)
	for _, f := range files {
		if f.Name() == filename {
			file, err = os.ReadFile(wd)
			if err != nil {
				log.Error(err)
				return nil, "", err
			}

			fname = f.Name()
		}
	}

	return
}

func (r *repo) InsertMediaData(data *Media) error {
	props := mongohelper.InsertOneProps{
		Ctx:            context.TODO(),
		CollectionName: metadataCollectionName,
		InsertData:     data,
	}

	_, err := r.dbHelper.InsertOne(props)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (r *repo) FindMediaData(filter primitive.D, opts options.FindOptions) (*mongo.Cursor, error) {
	props := mongohelper.FindProps{
		Ctx:            context.TODO(),
		CollectionName: metadataCollectionName,
		Filter:         filter,
		Opts:           opts,
	}

	res, err := r.dbHelper.Find(props)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return res, nil
}

func (r *repo) UpdateMediaData(filter primitive.D, mediaData *UpdateMedia) (*mongo.UpdateResult, error) {
	dataByte, err := bson.Marshal(mediaData)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var bsonData bson.M
	err = bson.Unmarshal(dataByte, &bsonData)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	props := mongohelper.UpdateOneProps{
		Ctx:            context.TODO(),
		CollectionName: metadataCollectionName,
		Filter:         filter,
		UpdateData:     bson.D{{Key: "$set", Value: bsonData}},
	}
	res, err := r.dbHelper.UpdateOne(props)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return res, nil
}

func (r *repo) DeleteMediaData(filter primitive.D) (*mongo.DeleteResult, error) {
	props := mongohelper.DeleteOneProps{
		Ctx:            context.TODO(),
		CollectionName: metadataCollectionName,
		Filter:         filter,
	}

	res, err := r.dbHelper.DeleteOne(props)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return res, nil
}

type repo struct {
	dbHelper mongohelper.MongoDBHelper
}

func NewRepository(helper mongohelper.MongoDBHelper) Repository {
	r := &repo{
		dbHelper: helper,
	}

	return r
}
