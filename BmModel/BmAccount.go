package BmModel

import (
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	Account  string `json:"account" bson:"account"`
	Password string `json:"password" bson:"password"`
	BrandId  string `json:"brand-id" bson:"brand-id"`
	Token    string `json:"token"`
}

func (u *Account) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
