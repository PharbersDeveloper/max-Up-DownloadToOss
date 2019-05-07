package BmDataStorage

import (
	"errors"
	"fmt"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons"
	"github.com/alfredyang1986/BmServiceDef/BmDaemons/BmMongodb"
	"github.com/manyminds/api2go"
	"net/http"
	"github.com/PharbersDeveloper/max-Up-DownloadToOss/BmModel"
)

// BmAccountStorage stores all of the tasty modelleaf, needs to be injected into
// Account and Account Resource. In the real world, you would use a database for that.
type BmAccountStorage struct {
	db *BmMongodb.BmMongodb
}

func (s BmAccountStorage) NewAccountStorage(args []BmDaemons.BmDaemon) *BmAccountStorage {
	mdb := args[0].(*BmMongodb.BmMongodb)
	return &BmAccountStorage{db: mdb}
}

// GetAll of the modelleaf
func (s BmAccountStorage) GetAll(r api2go.Request, skip int, take int) []BmModel.Account {
	in := BmModel.Account{}
	var out []BmModel.Account
	err := s.db.FindMulti(r, &in, &out, skip, take)
	if err == nil {
		for i, iter := range out {
			s.db.ResetIdWithId_(&iter)
			out[i] = iter
		}
		return out
	} else {
		return nil
	}
}

// GetOne tasty modelleaf
func (s BmAccountStorage) GetOne(id string) (BmModel.Account, error) {
	in := BmModel.Account{ID: id}
	out := BmModel.Account{ID: id}
	err := s.db.FindOne(&in, &out)
	if err == nil {
		return out, nil
	}
	errMessage := fmt.Sprintf("Account for id %s not found", id)
	return BmModel.Account{}, api2go.NewHTTPError(errors.New(errMessage), errMessage, http.StatusNotFound)
}

// Insert a fresh one
func (s *BmAccountStorage) Insert(c BmModel.Account) string {
	tmp, err := s.db.InsertBmObject(&c)
	if err != nil {
		fmt.Println(err)
	}

	return tmp
}

// Delete one :(
func (s *BmAccountStorage) Delete(id string) error {
	in := BmModel.Account{ID: id}
	err := s.db.Delete(&in)
	if err != nil {
		return fmt.Errorf("Account with id %s does not exist", id)
	}

	return nil
}

// Update updates an existing modelleaf
func (s *BmAccountStorage) Update(c BmModel.Account) error {
	err := s.db.Update(&c)
	if err != nil {
		return fmt.Errorf("Account with id does not exist")
	}

	return nil
}
