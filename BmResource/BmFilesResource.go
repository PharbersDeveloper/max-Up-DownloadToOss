package BmResource

import (
	"errors"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmDataStorage"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmModel"
	"github.com/manyminds/api2go"
	"net/http"
	"reflect"
)

type BmFilesResource struct {
	BmFilesStorage *BmDataStorage.BmFilesStorage
}

func (c BmFilesResource) NewFilesResource(args []BmDataStorage.BmStorage) BmFilesResource {
	var cs *BmDataStorage.BmFilesStorage
	for _, arg := range args {
		tp := reflect.ValueOf(arg).Elem().Type()
		if tp.Name() == "BmFilesStorage" {
			cs = arg.(*BmDataStorage.BmFilesStorage)
		}	
	}
	return BmFilesResource{BmFilesStorage: cs}
}

// FindAll Filess
func (c BmFilesResource) FindAll(r api2go.Request) (api2go.Responder, error) {
	result := c.BmFilesStorage.GetAll(r,-1,-1)
	return &Response{Res: result}, nil
}

// FindOne choc
func (c BmFilesResource) FindOne(ID string, r api2go.Request) (api2go.Responder, error) {
	res, err := c.BmFilesStorage.GetOne(ID)
	return &Response{Res: res}, err
}

// Create a new choc
func (c BmFilesResource) Create(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(BmModel.Files)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	id := c.BmFilesStorage.Insert(choc)
	choc.ID = id
	return &Response{Res: choc, Code: http.StatusCreated}, nil
}

// Delete a choc :(
func (c BmFilesResource) Delete(id string, r api2go.Request) (api2go.Responder, error) {
	err := c.BmFilesStorage.Delete(id)
	return &Response{Code: http.StatusOK}, err
}

// Update a choc
func (c BmFilesResource) Update(obj interface{}, r api2go.Request) (api2go.Responder, error) {
	choc, ok := obj.(BmModel.Files)
	if !ok {
		return &Response{}, api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	err := c.BmFilesStorage.Update(choc)
	return &Response{Res: choc, Code: http.StatusNoContent}, err
}
