package BmDataStorage

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmModel"
	"github.com/manyminds/api2go"
)

// BmFilesStorage stores all Fileses
type BmFilesStorage struct {
	db *BmMongodb.BmMongodb
}

func (s BmFilesStorage) NewFilesStorage(args []BmDaemons.BmDaemon) *BmFilesStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmFilesStorage{mdb}
}

// GetAll returns the model map (because we need the ID as key too)
func (s BmFilesStorage) GetAll(r api2go.Request, skip int, take int) []*BmModel.Files {
	in := BmModel.Files{}
	var out []BmModel.Files
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		var tmp []*BmModel.Files
		for i := 0; i < len(out); i++ {
			ptr := out[i]
			s.db.ResetIdWithId_(&ptr)
			tmp = append(tmp, &ptr)
		}
		return tmp
	} else {
		return nil //make(map[string]*BmModel.Files)
	}
}

// GetOne model
func (s BmFilesStorage) GetOne(id string) (BmModel.Files, error) {
	in := BmModel.Files{ID: id}
	model := BmModel.Files{ID: id}
	err := s.db.FindOne(&in, &model)
	if err == nil {
		return model, nil
	}
	errMessage := fmt.Sprintf("Files for id %s not found", id)
	return BmModel.Files{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a model
func (s *BmFilesStorage) Insert(c BmModel.Files) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmFilesStorage) Delete(id string) error {
	in := BmModel.Files{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Files with id %s does not exist", id)
	}

	return nil
}

// Update a model
func (s *BmFilesStorage) Update(c BmModel.Files) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Files with id does not exist")
	}

	return nil
}

func (s *BmFilesStorage) Count(req api2go.Request, c BmModel.Files) int {
	r, _ := s.db.Count(req, &c)
	return r
}
