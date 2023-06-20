package content

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"appota/web-builder/config"
	"appota/web-builder/db/mongohelper"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	templateCollectionName = "template"
)

type Repository interface {
	InsertTemplatesData(data *Template) error
	FindTemplatesData(filter primitive.D, opts options.FindOptions) (*mongo.Cursor, error)
	DeleteTemplatesData(id string) (*mongo.DeleteResult, error)
	SaveJSONFile(name string, file []byte, path string) error
	GetJSONFile(id string, path string) ([]byte, error)
	DeleteJSONFile(id string, path string) error
	UpdateTemplatesData(id string, data *Template) (*Template, error)
}

type repo struct {
	dbHelper mongohelper.MongoDBHelper
}

func (r *repo) InsertTemplatesData(data *Template) error {
	props := mongohelper.InsertOneProps{
		Ctx:            context.TODO(),
		CollectionName: templateCollectionName,
		InsertData:     data,
	}

	_, err := r.dbHelper.InsertOne(props)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (r *repo) FindTemplatesData(filter primitive.D, opts options.FindOptions) (*mongo.Cursor, error) {
	props := mongohelper.FindProps{
		Ctx:            context.TODO(),
		CollectionName: templateCollectionName,
		Filter:         filter,
		Opts:           opts,
	}

	res, err := r.dbHelper.Find(props)
	if err != nil {
		log.Error(err)
		return &mongo.Cursor{}, err
	}

	return res, nil
}

func (r *repo) DeleteTemplatesData(id string) (*mongo.DeleteResult, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	props := mongohelper.DeleteOneProps{
		Ctx:            context.TODO(),
		CollectionName: templateCollectionName,
		Filter:         bson.D{{Key: "_id", Value: objectID}},
	}

	res, err := r.dbHelper.DeleteOne(props)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return res, nil
}

func (r *repo) UpdateTemplatesData(id string, data *Template) (*Template, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	dataByte, err := bson.Marshal(data)
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
		CollectionName: templateCollectionName,
		Filter:         bson.D{{Key: "_id", Value: objectID}},
		UpdateData:     bson.D{{Key: "$set", Value: bsonData}},
	}
	_, err = r.dbHelper.UpdateOne(props)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	data.ID = objectID

	return data, nil
}

func (r *repo) SaveJSONFile(name string, file []byte, path string) error {
	p := filepath.Join(config.GetStorePath(), path)
	if err := os.MkdirAll(p, os.ModePerm); err != nil {
		return err
	}

	wd := filepath.Join(p, name)
	dist, err := os.Create(wd)
	if err != nil {
		log.Error(err)
		return err
	}
	defer dist.Close()

	fileSource := bytes.NewReader(file)
	if _, err := io.Copy(dist, fileSource); err != nil {
		return err
	}

	return nil
}

func (r *repo) GetJSONFile(id string, path string) (file []byte, err error) {
	filename := strings.Join([]string{id, "json"}, ".")
	wd := filepath.Join(config.GetStorePath(), path, filename)

	file, err = os.ReadFile(wd)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return file, nil
}

func (r *repo) DeleteJSONFile(id string, path string) error {
	filename := strings.Join([]string{id, "json"}, ".")
	wd := filepath.Join(config.GetStorePath(), path, filename)

	err := os.Remove(wd)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func NewRepository(helper mongohelper.MongoDBHelper) Repository {
	r := &repo{
		dbHelper: helper,
	}

	return r
}
