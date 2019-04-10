package BmModel

import (
	bson "gopkg.in/mgo.v2/bson"
	//"fmt"
)

// File is the File that a user consumes in order to get fat and happy
type Files struct {
	ID   string         `json:"-"`
	Id_  bson.ObjectId  `json:"-" bson:"_id"`
	Name  string        `json:"name" bson:"name"`
	UploadTime  string `json:"uploadtime" bson:"uploadtime"`
	Describe  string    `json:"describe" bson:"describe"`
	Accept  string    	`json:"accept" bson:"accept"`
	Uuid  string    	`json:"uuid" bson:"uuid"`
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (c Files) GetID() string {
	return c.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (c *Files) SetID(id string) error {
	c.ID = id
	return nil
}

func (u *Files) GetConditionsBsonM(parameters map[string][]string) bson.M {
	rst := make(map[string]interface{})
	for k, v := range parameters {
		switch k {
		case "accept":
			rst[k] = v[0]
		}
	}
	return rst
}
