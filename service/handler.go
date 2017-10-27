package service

import (
	"github.com/unrolled/render"
	"io/ioutil"
	"net/http"
)

func createFirmPersonHandler(formater *render.Render, repository userWhoRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		bytes, err := ioutil.ReadAll(req.Body)

	}
}
