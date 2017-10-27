package service

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func NewServer() *negroni.Negroni {

	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	reuter := mux.NewRouter()

	repository := newInMemoryUserWhoRepository

	initRouter(reuter, formatter, repository())

	n.UseHandler(reuter)

	return n
}

func initRouter(mux *mux.Router, render *render.Render, repository userWhoRepository) {

}
