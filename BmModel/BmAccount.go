package BmModel

import (
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	EmployeeID string `json:"employee-id" bson:"employee-id"`
	Email      string `json:"email" bson:"email"`
	Phone      string `json:"phone" bson:"phone"`
	Username   string `json:"username" bson:"username"`
	Password   string `json:"password" bson:"password"`
	Account    string `json:"account" bson:"account"`
	
	RegisterDate float64 `json:"register-date" bson:"register-date"`
}

func (u *Account) GetConditionsBsonM(parameters map[string][]string) bson.M {
	return bson.M{}
}
