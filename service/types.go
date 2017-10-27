package service

import "github.com/liambilbo/userwho-engine"

type newPhysicalPersonRequest struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	SecondSurname string `json:"secondsurname"`
}

type newPhysicalPersonResponse struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	SecondSurname string `json:"secondsurname"`
}

type newFirmPersonRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type newFirmPersonResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (preq *newFirmPersonResponse) copyPerson(person *userwho_engine.FirmPersonImpl) {
	preq.Id = person.Id
	preq.Name = person.Name
}

func (preq *newPhysicalPersonResponse) copyPerson(person *userwho_engine.PhysicalPersonImpl) {
	preq.Id = person.Id
	preq.Name = person.Name
	preq.Surname = person.Surname
	preq.SecondSurname = person.SecondSurname
}

type userWhoRepository interface {
	getPersonById(id string) (person userwho_engine.Person, err error)
	getPersons() (persons []userwho_engine.Person, err error)
	addPerson(person userwho_engine.Person) (err error)
	updatePerson(id string, person userwho_engine.Person) (err error)
}
