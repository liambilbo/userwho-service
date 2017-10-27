package service

import (
	"encoding/json"
	"github.com/liambilbo/userwho-engine"
	"github.com/unrolled/render"
	"io/ioutil"
	"net/http"
)

func createFirmPersonHandler(formatter *render.Render, repository userWhoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var new newFirmPersonRequest
		err := json.Unmarshal(payload, &new)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse request")
			return
		}

		newPerson := userwho_engine.NewFirmPerson(new.Name)
		err = repository.addPerson(newPerson)
		if err != nil {
			formatter.Text(w, http.StatusNotModified, "Failed to insert Firm Person")
			return
		}

		var newPersonResponse newFirmPersonResponse
		newPersonResponse.copyPerson(newPerson)
		w.Header().Add("Location", "/persons/"+newPersonResponse.Id)
		formatter.JSON(w, http.StatusCreated, newPersonResponse)
	}
}

func createPhysicalPersonHandler(formatter *render.Render, repository userWhoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var new newPhysicalPersonRequest
		err := json.Unmarshal(payload, &new)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse request")
			return
		}

		newPerson := userwho_engine.NewPhysicalPerson(new.Name, new.Surname, new.SecondSurname)
		err = repository.addPerson(newPerson)
		if err != nil {
			formatter.Text(w, http.StatusNotModified, "Failed to insert Firm Person")
			return
		}

		var newPersonResponse newPhysicalPersonResponse
		newPersonResponse.copyPerson(newPerson)
		w.Header().Add("Location", "/person/"+newPersonResponse.Id)
		formatter.JSON(w, http.StatusOK, newPersonResponse)
	}
}
