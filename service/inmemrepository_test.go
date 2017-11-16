package service

import (
	"github.com/liambilbo/userwho-engine"
	"testing"
	"time"
)

func TestGetPersonFromRepositoryById(t *testing.T) {

	repo := newInMemoryUserWhoRepository()

	issueDate := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	document := userwho_engine.NewDocument("K1234", "CIF", "ESP", &issueDate, nil)
	fiscalAdress := userwho_engine.NewAddress("ESP", "28029", "Madrid",
		"Madrid", "CL", "Gutierrez Ceciana", "23", "")
	firmPerson := userwho_engine.NewFirmPerson("Caminos S.A", "ESP", "ESP", document, fiscalAdress)
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
