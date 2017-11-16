package service

import (
	"github.com/cloudnativego/cfmgo"
	"github.com/liambilbo/userwho-engine"
	"github.com/liambilbo/userwho-service/fakes"
	"testing"
)

var (
	fakeDBURI = "mongodb://fake.uri@addr:port/guid"
)

func TestAddFirmPersonShowsUpInMongoRepository(t *testing.T) {
	var fakeFirms []firmRecord
	var userWhoCollection = cfmgo.Connect(
		fakes.FakeNewCollectionDialer(fakeFirms),
		fakeDBURI,
		PersonsCollectionName)

	repo := newMongoUserWhoRepository(userWhoCollection)

	address := userwho_engine.NewAddress("ESP", "28027", "Madrid", "Madrid", "CL", "Pintor Lokela", "27", "")
	document := userwho_engine.NewDocument("PK23456", "CIF", "ESP", nil, nil)
	firmPerson := userwho_engine.NewFirmPerson("Natural House S.A.", "ESP", "ESP", document, address)

	err := repo.addPerson(firmPerson)

	if err != nil {
		t.Errorf("Error adding firm to mongo: %v", err)
	}

	persons, err := repo.getPersons()
	if err != nil {
		t.Errorf("Got an error retrieving persons: %v", err)
	}
	if len(persons) != 1 {
		t.Errorf("Expected persons length to be 1; received %d", len(persons))
	}
}
