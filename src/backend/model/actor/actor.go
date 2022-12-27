package actor

import (
	"github.com/lpuig/batec/stockmanagement/src/backend/model/date"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/timestamp"
)

type Actor struct {
	Id        int
	Ref       string
	FirstName string
	LastName  string
	State     string
	Period    date.DateStringRange
	Company   string
	Comment   string
	timestamp.TimeStamp
}

func NewActor(firstName, lastName, company string) *Actor {
	return &Actor{
		Id:        -1,
		Ref:       lastName + " " + firstName,
		FirstName: firstName,
		LastName:  lastName,
		State:     "",
		Period:    date.DateStringRange{},
		Company:   company,
		Comment:   "",
	}
}
