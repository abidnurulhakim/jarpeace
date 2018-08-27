package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Model struct {
	ID        bson.ObjectId `mgorm:"primary_key" bson:"_id,omitempty" json:"id"`
	UpdatedAt time.Time     `bson:"updated_at,omitempty" json:"updated_at"`
	CreatedAt time.Time     `mgorm:"index" bson:"created_at,omitempty" json:"created_at"`
}
