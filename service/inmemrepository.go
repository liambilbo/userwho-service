package service

import (
	"errors"
	"github.com/liambilbo/userwho/userwho-engine"
)

type inMemoryUserWhoRepository struct {
	persons []userwho_engine.Person
}

func newInMemoryUserWhoRepository() *inMemoryUserWhoRepository {
	repository := inMemoryUserWhoRepository{}
	repository.persons = []userwho_engine.Person{}
	return &repository
}

func (repository *inMemoryUserWhoRepository) getPersonById(id string) (person userwho_engine.Person, err error) {
	for _, v := range repository.persons {
		if v.GetId() == id {
			person = v
			return
		}
	}
	err = errors.New("could not found in repository")
	return
}

func (repository *inMemoryUserWhoRepository) getPersons() (persons []userwho_engine.Person, err error) {
	persons = repository.persons
	return
}

func (repository *inMemoryUserWhoRepository) addPerson(person userwho_engine.Person) (err error) {
	found := false
	for _, v := range repository.persons {
		if v.GetId() == person.GetId() {
			found = true
			break
		}
	}

	if !found {
		repository.persons = append(repository.persons, person)
	} else {
		err = errors.New("Person already exists")
	}
	return

}

func (repository *inMemoryUserWhoRepository) updatePerson(id string, person userwho_engine.Person) (err error) {

	found := false
	for k, v := range repository.persons {
		if v.GetId() == person.GetId() {
			found = true
			repository.persons[k] = person
			break
		}
	}

	if !found {
		err = errors.New("Person not found")
	}
	return

}
