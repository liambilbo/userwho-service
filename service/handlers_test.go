package service

import (
	"bytes"
	"encoding/json"
	"github.com/unrolled/render"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

const (
	fakeMatchLocationResult = "/firmpersons/5a003b78-409e-4452-b456-a6f0dcee05bd"
)

func CreateFirmPersonRespondToBadData(t *testing.T) {
	client := &http.Client{}
	repository := newInMemoryUserWhoRepository()
	server := httptest.NewServer(http.HandlerFunc(createFirmPersonHandler(formatter, repository)))
	defer server.Close()

	body1 := []byte("this is not valid json")
	body2 := []byte("{\"test\":\"this is valid json, but doesn't conform to server expectations.\"}")

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body1))
	if err != nil {
		t.Errorf("Error in creating POST request for createMatchHandler: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to createMatchHandler: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Error("Sending invalid JSON should result in a bad request from server.")
	}

	req2, err2 := http.NewRequest("POST", server.URL, bytes.NewBuffer(body2))
	if err2 != nil {
		t.Errorf("Error in creating second POST request for invalid data on create match: %v", err2)
	}
	req2.Header.Add("Content-Type", "application/json")
	res2, _ := client.Do(req2)
	defer res2.Body.Close()
	if res2.StatusCode != http.StatusBadRequest {
		t.Error("Sending valid JSON but with incorrect or missing fields should result in a bad request and didn't.")
	}
}

func TestCreateFirmPerson(t *testing.T) {
	client := &http.Client{}
	person := []byte("{\"name\":\"Lorea Gardening S.A.\"}")
	repository := newInMemoryUserWhoRepository()
	server := httptest.NewServer(createFirmPersonHandler(formatter, repository))

	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(person))
	if err != nil {
		t.Errorf("Error in creating request for firm person %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in sending request for firm person %v", err)
	}
	defer res.Body.Close()

	payload, _ := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusCreated {
		t.Errorf("Error creating firm expected code 201 , instead received %s ", res.StatusCode)
	}

	loc, headerOk := res.Header["Location"]

	if !headerOk {
		t.Error("Location header is not set")
	} else {
		if !strings.Contains(loc[0], "/firmpersons/") {
			t.Errorf("Location header should contain '/persons/'")
		}
		if len(loc[0]) != len(fakeMatchLocationResult) {
			t.Errorf("Location value does not contain guid of new person")
		}
	}

	var newPerson newFirmPersonResponse
	err = json.Unmarshal(payload, &newPerson)
	if err != nil {
		t.Errorf("Error Unmarsalling Firm %v ", err)
	}

	if newPerson.Name != "Lorea Gardening S.A." {
		t.Errorf("Error Unmarsalling Firm %v ", err)
	}

	if newPerson.Id == "" || !strings.Contains(loc[0], newPerson.Id) {
		t.Error("newPerson.Id does not match Location header")
	}

}
