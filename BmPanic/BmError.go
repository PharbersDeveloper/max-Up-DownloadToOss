package BmPanic

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/go-uuid"
	"github.com/manyminds/api2go"
	"net/http"
	"strconv"
	"sync"
)

var ALFRED_TEST_ERROR = errors.New("alfred test error")

type tBMError struct {
	m map[string]api2go.HTTPError
}

var e *tBMError
var o sync.Once

//TODO:error definition load from file
func ErrInstance() *tBMError {
	o.Do(func() {
		e = &tBMError{
			m: map[string]api2go.HTTPError{
				"Auth Failed!": api2go.HTTPError{Errors: []api2go.Error{
					{
						Links: &api2go.ErrorLinks{
							About: "http://login",
						},
						Status: "401",
						Code:   "001",
						Title:  "Auth error!",
						Detail: "No token found!",
						Source: &api2go.ErrorSource{
							Pointer: "#titleField",
						},
						Meta: map[string]interface{}{
							"creator": "jeorch",
						},
					},
				}},
				"token not exist": api2go.HTTPError{Errors: []api2go.Error{
					{
						Links: &api2go.ErrorLinks{
							About: "http://login",
						},
						Status: "401",
						Code:   "001",
						Title:  "Auth error!",
						Detail: "token not exist!",
						Source: &api2go.ErrorSource{
							Pointer: "#titleField",
						},
						Meta: map[string]interface{}{
							"creator": "jeorch",
						},
					},
				}},
				"no defind error!": api2go.HTTPError{Errors: []api2go.Error{
					{
						Links: &api2go.ErrorLinks{
							About: "http://404",
						},
						Status: "400",
						Code:   "9999",
						Title:  "no defind error!",
						Detail: "no defind error!",
						Source: &api2go.ErrorSource{
							Pointer: "#titleField",
						},
						Meta: map[string]interface{}{
							"creator": "jeorch",
						},
					},
				}},
			},
		}
	})
	return e
}

func (e *tBMError) IsErrorDefined(ec string) bool {
	for k, _ := range e.m {
		if k == ec {
			return true
		}
	}
	return false
}

func resetlHTTPErrorID(input api2go.HTTPError) {
	if len(input.Errors) != 0 {
		for i,_ := range input.Errors {
			eid, _ := uuid.GenerateUUID()
			input.Errors[i].ID = eid
		}
	}
}

func (e *tBMError) ErrorReval(err interface{}, w http.ResponseWriter) {
	es := err.(string)
	var hr api2go.HTTPError
	if e.IsErrorDefined(es) {
		hr = e.m[es]
	} else {
		hr = e.m["no defind error!"]
		hr.Errors[0].Detail = es
	}
	resetlHTTPErrorID(hr)
	enc := json.NewEncoder(w)
	w.Header().Add("Content-Type", "application/json")
	statusCode,  _ := strconv.Atoi(hr.Errors[0].Status)
	w.WriteHeader(statusCode)
	enc.Encode(hr)
}