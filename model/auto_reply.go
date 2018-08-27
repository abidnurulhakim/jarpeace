package model

type AutoReply struct {
	Model
	Text   string `bson:"text" json:"text"`
	Answer string `bson:"answer" json:"answer"`
}

func NewAutoReply() AutoReply {
	autoReply := AutoReply{}
	return autoReply
}
