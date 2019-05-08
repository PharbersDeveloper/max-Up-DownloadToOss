package BmResource

import (
	"errors"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmModel"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmDataStorage"
)

type BmAccountResource struct {
	BmAccountStorage *BmDataStorage.BmAccountStorage
}

func (c BmAccountResource) NewAccountResource(args []BmDataStorage.BmStorage) *BmAccountResource {
	var cs *BmDataStorage.BmAccountStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmAccountStorage" {
			cs = arg.(*BmDataStorage.BmAccountStorage)
		}
	}
	return &BmAccountResource{
		BmAccountStorage: cs,
	}
}

// FindAll images
func (c BmAccountResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	var result []BmModel.Account
	result = c.BmAccountStorage.GetAll(r, -1, -1)
	return &Response{Res: result}, nil
}

// FindOne account
func (c BmAccountResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.BmAccountStorage.GetOne(ID)
	return &Response{Res: res}, err
}

// Create a new account
func (c BmAccountResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	account, ok := obj.(BmModel.Account)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest,
		)
	}

	id := c.BmAccountStorage.Insert(account)
	account.ID = id
	return &Response{Res: account, Code: http.StatusCreated}, nil
}

// Delete a account :(
func (c BmAccountResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.BmAccountStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a account
func (c BmAccountResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	account, ok := obj.(BmModel.Account)
	if !ok {
		return &Response{}, api2go.NewHTTPError(
			errors.New("Invalid instance given"),
			"Invalid instance given",
			http.StatusBadRequest,
		)
	}

	err := c.BmAccountStorage.Update(account)
	return &Response{Res: account, Code: http.StatusNoContent}, err
}
