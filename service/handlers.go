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
		var newf newFirmPersonRequest
		err := json.Unmarshal(payload, &newf)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse request")
			return
		}

		document, address := convertToEngine(newf.newPersonRequest)

		newPerson := userwho_engine.NewFirmPerson(newf.Name,
			userwho_engine.Country(newf.Nationality),
			userwho_engine.Country(newf.SettingUpCountry),
			document, address)

		err = repository.addPerson(newPerson)
		if err != nil {
			formatter.Text(w, http.StatusNotModified, "Failed to insert Firm Person")
			return
		}

		var newPersonResponse newFirmPersonResponse
		newPersonResponse.copyPerson(newPerson)
		w.Header().Add("Location", "/firmpersons/"+newPersonResponse.Id)
		formatter.JSON(w, http.StatusCreated, newPersonResponse)
	}
}

func createPhysicalPersonHandler(formatter *render.Render, repository userWhoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		payload, _ := ioutil.ReadAll(req.Body)
		var newP newPhysicalPersonRequest
		err := json.Unmarshal(payload, &newP)
		if err != nil {
			formatter.Text(w, http.StatusBadRequest, "Failed to parse request")
			return
		}

		document, address := convertToEngine(newP.newPersonRequest)
		newPerson := userwho_engine.NewPhysicalPerson(newP.Name, newP.Surname, newP.SecondSurname,
			userwho_engine.Country(newP.Nationality), userwho_engine.Country(newP.Address.Country), document, address)
		err = repository.addPerson(newPerson)
		if err != nil {
			formatter.Text(w, http.StatusNotModified, "Failed to insert Firm Person")
			return
		}

		var newPersonResponse newPhysicalPersonResponse
		newPersonResponse.copyPerson(newPerson)
		w.Header().Add("Location", "/physicalpersons/"+newPersonResponse.Id)
		formatter.JSON(w, http.StatusOK, newPersonResponse)
	}
}

func convertToEngine(new newPersonRequest) (document userwho_engine.Document, address userwho_engine.Address) {

	document = userwho_engine.NewDocument(new.Document.Number,
		userwho_engine.DocumentType(new.Document.Type),
		userwho_engine.Country(new.Document.IssueCountry),
		&new.Document.IssueDate, &new.Document.MaturityDate)

	address = userwho_engine.NewAddress(userwho_engine.Country(new.Address.Country),
		new.Address.PostalCode,
		new.Address.Province,
		new.Address.Town,
		userwho_engine.StreetType(new.Address.StreetType),
		new.Address.Street,
		new.Address.StreetNumber,
		new.Address.Complementary)
	return
}
