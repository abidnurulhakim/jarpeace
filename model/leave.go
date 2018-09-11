package model

import (
	"strings"

	"github.com/abidnurulhakim/carbon"
	"github.com/abidnurulhakim/jarpeace/helper"
	"gopkg.in/mgo.v2/bson"
)

type Leave struct {
	ID        bson.ObjectId  `bson:"_id,omitempty" json:"id"`
	Type      string         `bson:"type" json:"type"`
	Start     *carbon.Carbon `bson:"start" json:"start"`
	End       *carbon.Carbon `bson:"end" json:"end"`
	Note      string         `bson:"note" json:"note"`
	UserID    int            `bson:"user_id" json:"user_id"`
	Username  string         `bson:"username" json:"username"`
	UpdatedAt *carbon.Carbon `bson:"updated_at,omitempty" json:"updated_at"`
	CreatedAt *carbon.Carbon `bson:"created_at,omitempty" json:"created_at"`
	DeletedAt *carbon.Carbon `bson:"deleted_at,omitempty" json:"deleted_at"`
}

func NewLeave(raw string) Leave {
	leave := Leave{}
	if raw == "" {
		return leave
	}
	now := carbon.Now()
	params := strings.Split(raw, ";")
	if len(params) > 0 {
		leave.Type = strings.ToLower(strings.TrimSpace(params[0]))
	}
	if len(params) > 1 {
		startDate, _ := helper.ParseHumanDatetime(params[1], now.StartOfDay())
		leave.Start = startDate.AddSeconds(1)
	}
	if len(params) > 2 {
		endDate, _ := helper.ParseHumanDatetime(params[2], now.EndOfDay())
		leave.End = endDate
	}
	if len(params) > 3 {
		leave.Note = strings.TrimSpace(params[3])
	}
	return leave
}
