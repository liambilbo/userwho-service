package service

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/cloudnativego/cf-tools"
	"github.com/cloudnativego/cfmgo"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

// NewServer configures and returns a Server.
func NewServer(appEnv *cfenv.App) *negroni.Negroni {

	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	n := negroni.Classic()
	mx := mux.NewRouter()

	repo := initRepository(appEnv)

	initRoutes(mx, formatter, repo)

	n.UseHandler(mx)
	return n
}

func initRepository(appEnv *cfenv.App) (repo userWhoRepository) {
	dbServiceURI, err := cftools.GetVCAPServiceProperty(dbServiceName, "url", appEnv)
	if err != nil || dbServiceURI == "" {
		if err != nil {
			fmt.Printf("\nError retrieving database configuration: %v\n", err)
		}
		fmt.Println("MongoDB was not detected; configuring inMemoryRepository...")
		repo = newInMemoryUserWhoRepository()
		return
	}
	matchCollection := cfmgo.Connect(cfmgo.NewCollectionDialer, dbServiceURI, PersonsCollectionName)
	fmt.Printf("Connecting to MongoDB service: %s...\n", dbServiceName)
	repo = newMongoUserWhoRepository(matchCollection)
	return
}

func initRoutes(mux *mux.Router, formatter *render.Render, repository userWhoRepository) {
	mux.HandleFunc("/firmpersons", createFirmPersonHandler(formatter, repository)).Methods("POST")
	mux.HandleFunc("/physicalpersons", createFirmPersonHandler(formatter, repository)).Methods("POST")
}
