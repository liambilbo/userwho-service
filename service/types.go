package service

import "github.com/liambilbo/userwho-engine"

type newMatchResponse struct {
	ID          string `json:"id"`
	StartedAt   int64  `json:"started_at"`
	GridSize    int    `json:"gridsize"`
	PlayerWhite string `json:"playerWhite"`
	PlayerBlack string `json:"playerBlack"`
	Turn        int    `json:"turn,omitempty"`
}

type userWhoRepository interface {
	getPersonById(id string) (person userwho_engine.Person, err error)
	getPersons() (persons []userwho_engine.Person, err error)
	addPerson(person userwho_engine.Person) (err error)
	updatePerson(id string, person userwho_engine.Person) (err error)
}
