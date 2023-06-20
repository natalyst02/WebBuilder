package content

import "go.mongodb.org/mongo-driver/bson/primitive"

type Content struct {
	ID string
}
type Template struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      string             `bson:"type,omitempty" json:"type"`
	Name      string             `bson:"name,omitempty" json:"name"`
	Tags      string             `bson:"tags,omitempty" json:"tags"`
	ProjectID string             `bson:"projectId,omitempty" json:"projectId"`
	Content   interface{}        `bson:"-" json:"-"`
}
