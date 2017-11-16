package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudnativego/cfmgo"
	"github.com/cloudnativego/cfmgo/params"
	"github.com/liambilbo/userwho-engine"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type mongoUserWhoRepository struct {
	Collection cfmgo.Collection
}

type documentRecord struct {
	Number       string     `bson:"number" json:"number"`
	Type         string     `bson:"type" json:"type"`
	IssueCountry string     `bson:"issue_country" json:"issue_country"`
	IssueDate    *time.Time `bson:"issue_date,omitempty" json:"issue_date"`
	MaturityDate *time.Time `bson:"maturity_date,omitempty" json:"maturity_date"`
}

type addressRecord struct {
	Type    string `bson:"type" json:"type"`
	Order   int    `bson:"order" json:"order"`
	Address struct {
		Country       string `bson:"country" json:"country"`
		PostalCode    string `bson:"postalcode" json:"postalcode"`
		Province      string `bson:"province" json:"province"`
		Town          string `bson:"town" json:"town"`
		StreetType    string `bson:"street_type" json:"street_type"`
		Street        string `bson:"street" json:"street"`
		StreetNumber  string `bson:"streetnumber" json:"streetnumber"`
		Complementary string `bson:"complementary" json:"complementary"`
	}
}

type actorRecord struct {
	*firmRecord
	*physicalRecord
}

type personRecord struct {
	RecordID   bson.ObjectId             `bson:"_id,omitempty" json:"id"`
	PersonID   string                    `bson:"person_id",json:"person_id"`
	PersonType string                    `bson:"person_type",json:"person_type"`
	Documents  map[string]documentRecord `bson:"documents",json:"documents"`
	Addresses  []addressRecord           `bson:"addresses",json:"addresses"`
}

func (p personRecord) GetPersonType() string {
	return p.PersonType
}

type firmRecord struct {
	personRecord `bson:",inline""`
	Name         string `bson:"name",json:"name"`
}

type physicalRecord struct {
	personRecord  `bson:",inline""`
	Name          string `bson:"name",json:"name"`
	Surname       string `bson:"surname",json:"surname"`
	SecondSurname string `bson:"secondsurname",json:"secondsurname"`
}

func newMongoUserWhoRepository(col cfmgo.Collection) (repo *mongoUserWhoRepository) {
	repo = &mongoUserWhoRepository{
		Collection: col,
	}
	return
}

func convertDocumentToDocumentRecord(document userwho_engine.Document) (documentR documentRecord, err error) {
	documentR = documentRecord{}
	documentR.Number = document.Number
	documentR.Type = string(document.Type)
	documentR.Type = string(document.IssueCountry)
	documentR.IssueDate = document.IssueDate
	documentR.MaturityDate = document.MaturityDate
	return
}

func convertAddressToAddressRecord(address userwho_engine.PersonAddress) (addressR addressRecord, err error) {
	addressR = addressRecord{}
	addressR.Type = string(address.Type)
	addressR.Order = address.Order
	addressR.Address.Country = string(address.Country)
	addressR.Address.Province = address.Province
	addressR.Address.PostalCode = address.PostalCode
	addressR.Address.StreetType = string(address.StreetType)
	addressR.Address.Town = address.Town
	addressR.Address.Street = address.Street
	addressR.Address.StreetNumber = address.StreetNumber
	addressR.Address.Complementary = address.Complementary
	return
}

func convertDocumentRecordToDocument(documentR documentRecord) (document userwho_engine.Document, err error) {
	document = userwho_engine.Document{}
	document.Number = documentR.Number
	document.Type = userwho_engine.DocumentType(documentR.Type)
	document.IssueDate = documentR.IssueDate
	document.MaturityDate = documentR.MaturityDate
	return
}

func convertAddressRecordToAddress(addressR addressRecord) (address userwho_engine.PersonAddress, err error) {
	address = userwho_engine.PersonAddress{}
	address.Type = userwho_engine.AddressType(addressR.Type)
	address.Order = addressR.Order
	address.Country = userwho_engine.Country(addressR.Address.Country)
	address.Province = addressR.Address.Province
	address.PostalCode = addressR.Address.PostalCode
	address.StreetType = userwho_engine.StreetType(addressR.Address.StreetType)
	address.Town = addressR.Address.Town
	address.Street = addressR.Address.Street
	address.StreetNumber = addressR.Address.StreetNumber
	address.Complementary = addressR.Address.Complementary
	return
}

func convertPersonToPersonRecord(m interface{}) (mr interface{}, err error) {

	if f, ok := isFirmPerson(m); ok {
		firmR := &firmRecord{}
		firmR.RecordID = bson.NewObjectId()
		firmR.PersonID = f.Id
		firmR.PersonType = string(f.Type)
		firmR.Name = f.Name
		loadFromPersonDocuments(f.Person, &firmR.personRecord)
		loadFromPersonAdresses(f.Person, &firmR.personRecord)
		mr = firmR
	} else if p, ok := isPhysicalPerson(m); ok {
		phR := &physicalRecord{}
		phR.RecordID = bson.NewObjectId()
		phR.Name = p.Name
		phR.Surname = p.Surname
		phR.SecondSurname = p.SecondSurname
		loadFromPersonDocuments(p.Person, &phR.personRecord)
		loadFromPersonAdresses(p.Person, &phR.personRecord)
		mr = phR
	} else {
		err = fmt.Errorf("Unsupported type %T. Ignoring.\n", m)
	}

	return
}

func loadFromPersonDocuments(person userwho_engine.Person, personR *personRecord) {
	docs := make(map[string]documentRecord)
	for k, v := range person.Documents.Documents {
		docs[string(k)], _ = convertDocumentToDocumentRecord(v)
	}

	personR.Documents = docs
}
func loadFromRecordDocuments(personR personRecord, person userwho_engine.Person) {
	for k, v := range personR.Documents {
		person.Documents.Documents[userwho_engine.DocumentType(k)], _ = convertDocumentRecordToDocument(v)
	}
}

func loadFromPersonAdresses(person userwho_engine.Person, personR *personRecord) {
	for _, v := range person.Addresses.Addresses {
		addR, _ := convertAddressToAddressRecord(v)
		personR.Addresses = append(personR.Addresses, addR)
	}
}
func loadFromRecordAddresses(personR personRecord, person userwho_engine.Person) {
	for _, v := range personR.Addresses {
		add, _ := convertAddressRecordToAddress(v)
		person.Addresses.Addresses = append(person.Addresses.Addresses, add)
	}
}

func convertPersonRecordToPerson(m interface{}) (mr interface{}, err error) {
	if firmR, okf := isFirmRecord(m); okf {
		firm := userwho_engine.FirmPerson{}
		firm.Name = firmR.Name
		mr = firm
	} else if physicalR, okp := isPhysicalRecord(m); okp {
		ph := userwho_engine.PhysicalPerson{}
		ph.Name = physicalR.Name
		ph.Surname = physicalR.Surname
		ph.SecondSurname = ph.SecondSurname
		mr = ph
	} else {
		err = fmt.Errorf("Unsupported type %T. Ignoring.\n", m)
	}

	return
}

func (repository *mongoUserWhoRepository) getPerson(id string) (person userwho_engine.Actor, err error) {
	repository.Collection.Wake()
	personR, err := repository.getMongoPerson(id)
	if err == nil {
		actor, err := convertPersonRecordToPerson(personR)
		if err != nil {
			person = actor.(userwho_engine.Actor)
		}
	}
	return
}

func (repository *mongoUserWhoRepository) addPerson(person userwho_engine.Actor) (err error) {
	repository.Collection.Wake()
	actor, _ := convertPersonToPersonRecord(person)

	if f, ok := isFirmRecord(actor); ok {
		_, err = repository.Collection.UpsertID(f.RecordID, f)
	} else if p, ok := isPhysicalRecord(actor); ok {
		_, err = repository.Collection.UpsertID(p.RecordID, p)
	} else {
		err = errors.New(fmt.Sprintf("Unsupported type %T. Ignoring.\n", actor))
	}
	return
}

func (repository *mongoUserWhoRepository) updatePerson(id string, person userwho_engine.Actor) (err error) {
	repository.Collection.Wake()

	foundActor, err := repository.getMongoPerson(id)
	if err == nil {
		actor, _ := convertPersonToPersonRecord(person)

		if f, okf := isFirmRecord(actor); okf {
			f.RecordID = foundActor.RecordID
			_, err = repository.Collection.UpsertID(foundActor.RecordID, f)
		} else if p, okp := isPhysicalRecord(actor); okp {
			p.RecordID = foundActor.RecordID
			_, err = repository.Collection.UpsertID(foundActor.RecordID, p)
		} else {
			err = errors.New(fmt.Sprintf("Unsupported type %T. Ignoring.\n", actor))
		}
		return
	}
	return

}

func (r *mongoUserWhoRepository) getPersons() (persons []userwho_engine.Actor, err error) {
	r.Collection.Wake()
	var pr []actorRecord
	params := &params.RequestParams{}
	_, err = r.Collection.Find(params, &pr)
	if err == nil {
		persons = make([]userwho_engine.Actor, len(pr))
		for k, v := range pr {
			a, _ := convertPersonRecordToPerson(v)
			persons[k] = a.(userwho_engine.Actor)
		}
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

func isFirmRecord(a interface{}) (result firmRecord, ok bool) {
	var p *firmRecord
	var s actorRecord
	p, ok = a.(*firmRecord)
	if ok {
		result = *p
	} else {
		s, ok = a.(actorRecord)

		if ok && s.firmRecord != nil {
			result = *s.firmRecord
		}

	}

	return
}

func isPhysicalRecord(a interface{}) (result physicalRecord, ok bool) {
	var p *physicalRecord
	var s actorRecord
	p, ok = a.(*physicalRecord)
	if ok {
		result = *p
	} else {
		s, ok = a.(actorRecord)

		if ok && s.physicalRecord != nil {
			result = *s.physicalRecord
		}
	}

	return
}

func isFirmPerson(a interface{}) (result userwho_engine.FirmPerson, ok bool) {
	p, ok := a.(userwho_engine.FirmPerson)
	if ok {
		result = p
	}
	return
}

func isPhysicalPerson(a interface{}) (result userwho_engine.PhysicalPerson, ok bool) {
	p, ok := a.(userwho_engine.PhysicalPerson)
	if ok {
		result = p
	}
	return
}

func (pr *actorRecord) UnmarshalJSON(b []byte) error {
	var p personRecord
	if err := json.Unmarshal(b, &p); err != nil {
		return err
	}
	switch p.PersonType {
	case "physical":
		var ph physicalRecord
		if err := json.Unmarshal(b, &ph); err != nil {
			return err
		}
		pr.physicalRecord = &ph
	case "firm":
		var f firmRecord
		if err := json.Unmarshal(b, &f); err != nil {
			return err
		}
		pr.firmRecord = &f
	}

	return nil
}
