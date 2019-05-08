package BmModel

import (
	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	ID  string        `json:"-"`
	Id_ bson.ObjectId `json:"-" bson:"_id"`

	Email      string `json:"email" bson:"email"`
	Phone      string `json:"phone" bson:"phone"`
	Password   string `json:"password" bson:"password"`
	Account    string `json:"account" bson:"account"`

}
func (c Account) GetID() string {
	return c.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (c *Account) SetID(id string) error {
	c.ID = id
	return nil
}
func (u *Account) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	for k, v := range parameters {
		switch k {
		case "password":
			rst[k] = v[0]
		case "email":
			rst[k] = v[0]
		case "phone":
			rst[k] = v[0]
		}
	}
	return rst
}
