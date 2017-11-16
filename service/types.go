package service

import (
	"github.com/liambilbo/userwho-engine"
	"time"
)

type newAddress struct {
	Country       string `json:"country"`
	PostalCode    string `json:"postalcode"`
	Province      string `json:"province"`
	Town          string `json:"town"`
	StreetType    string `json:"streetype"`
	Street        string `json:"street"`
	StreetNumber  string `json:"streetnumber"`
	Complementary string `json:"complementary"`
}

type newDocument struct {
	Number       string    `json:"number"`
	Type         string    `json:"type"`
	IssueCountry string    `json:"issuecountry"`
	IssueDate    time.Time `json:"issuedate"`
	MaturityDate time.Time `json:"maturitydate"`
}

type newPersonRequest struct {
	Id          string      `json:"id"`
	Nationality string      `json:"nationality"`
	Document    newDocument `json:"document"`
	Address     newAddress  `json:"fiscaladdress"`
}

type newPhysicalPersonRequest struct {
	newPersonRequest
	Name                 string    `json:"name"`
	Surname              string    `json:"surname"`
	SecondSurname        string    `json:"secondsurname"`
	BirthCity            string    `json:"birthcity"`
	BirthDate            time.Time `json:"birthdate"`
	BirthCountry         string    `json:"birthcountry"`
	SecondNationality    string    `json:"secondnationality"`
	Sex                  string    `json:"sex"`
	MaritalStatus        string    `json:"maritalstatus"`
	EducationLevel       string    `json:"educationlevel"`
	ProfessionalActivity string    `json:"professionalactivity"`
	LaboralSituation     string    `json:"laboralsituation"`
}

type newPhysicalPersonResponse struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	SecondSurname string `json:"secondsurname"`
}

type newFirmPersonRequest struct {
	newPersonRequest
	Name             string    `json:"name"`
	SettingUpDate    time.Time `json:"settingupdate"`
	SettingUpCountry string    `json:"settingupcountry"`
	Cnae             string    `json:"cnae"`
}

type newFirmPersonResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (preq *newFirmPersonResponse) copyPerson(person userwho_engine.FirmPerson) {
	preq.Id = person.Id
	preq.Name = person.Name
}

func (preq *newPhysicalPersonResponse) copyPerson(person userwho_engine.PhysicalPerson) {
	preq.Id = person.Id
	preq.Name = person.Name
	preq.Surname = person.Surname
	preq.SecondSurname = person.SecondSurname
}

type userWhoRepository interface {
	getPerson(id string) (person userwho_engine.Actor, err error)
	getPersons() (persons []userwho_engine.Actor, err error)
	addPerson(person userwho_engine.Actor) (err error)
	updatePerson(id string, person userwho_engine.Actor) (err error)
}
