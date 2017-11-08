package service

import (
	"errors"
	"github.com/cloudnativego/cfmgo"
	"github.com/cloudnativego/cfmgo/params"
	"github.com/liambilbo/userwho-engine"
	"gopkg.in/mgo.v2/bson"
)

type mongoUserWhoRepository struct {
	Collection cfmgo.Collection
}

type personRecord struct {
	RecordID   bson.ObjectId `bson:"_id,omitempty" json:"id"`
	PersonID   string        `bson:"person_id",json:"person_id"`
	PersonType string        `bson:"person_type",json:"person_type"`
	Name       string        `bson:"name",json:"name"`
}

type firmRecord struct {
	personRecord `bson:",inline""`
}

type physicalRecord struct {
	personRecord  `bson:",inline""`
	Surname       string `bson:"surname",json:"surname"`
	SecondSurname string `bson:"secondsurname",json:"secondsurname"`
}

func convertPersonToPersonRecord(m interface{}) (mr interface{}, err error) {

	switch m.(type) {
	case userwho_engine.FirmPerson:
		firm := m.(userwho_engine.FirmPerson)
		fir := &firmRecord{}
		fir.Name = firm.Name
		mr = fir

	case userwho_engine.PhysicalPerson:
		physical := m.(userwho_engine.PhysicalPerson)
		phr := &physicalRecord{}
		phr.Name = physical.Name
		phr.Surname = physical.Surname
		phr.SecondSurname = phr.SecondSurname
		mr = phr

	default:
		err = errors.New("Type is not FirmPerson or PhysicalPerson")

	}

	return
}

func (r mongoUserWhoRepository) getMongoPerson(id string) (person personRecord, err error) {
	var persons []personRecord
	query := bson.M{"person_id": id}
	params := &params.RequestParams{
		Q: query,
	}

	count, err := r.Collection.Find(params, &persons)

	if count == 0 {
		err = errors.New("Person not found")
	}

	if err != nil {
		person = persons[0]
	}
	return

}
