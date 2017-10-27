package service

import (
	"github.com/liambilbo/userwho/userwho-engine"
	"testing"
)

func TestGetPersonFromRepositoryById(t *testing.T) {

	repo := newInMemoryUserWhoRepository()
	firmPerson := userwho_engine.NewFirmPerson("Caminos S.A")
	err := repo.addPerson(firmPerson)

	if err != nil {
		t.Errorf("Error in addPerson")
	}

	persons, err := repo.getPersons()

	if err != nil {
		t.Errorf("Unexpected error in getPersons : %s ", err.Error())
	}

	if len(persons) != 1 {
		t.Errorf("Expected to have 1 Person and it retrieve %d", len(persons))
	}

}
