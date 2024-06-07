package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ToDoList struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Task     string             `json:"task,omitempty" bson:"task,omitempty"`
	Complete bool               `json:"complete,omitempty" bson:"complete,omitempty"`
}
