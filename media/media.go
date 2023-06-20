package media

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Media struct {
	File        []byte             `bson:"-" json:"-"`
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Title       string             `bson:"title,omitempty" json:"title"`
	Filename    string             `bson:"filename,omitempty" json:"filename"`
	Filepath    string             `bson:"filepath,omitempty" json:"filepath"`
	Description string             `bson:"description,omitempty" json:"description"`
	Tags        string             `bson:"tags,omitempty" json:"tags"`
	ProjectID   string             `bson:"projectId,omitempty" json:"projectId"`
	UploadedAt  time.Time          `bson:"uploadedAt" json:"uploadedAt"`
}

type UpdateMedia struct {
	Title       string `bson:"title,omitempty" json:"title"`
	Description string `bson:"description,omitempty" json:"description"`
	Tags        string `bson:"tags,omitempty" json:"tags"`
	ProjectID   string `bson:"projectId,omitempty"  json:"projectId"`
}
