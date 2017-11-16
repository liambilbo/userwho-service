package main

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/liambilbo/userwho-service/service"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	appEnv, err := cfenv.Current()
	if err != nil {
		fmt.Println("CF Environment not detected.")
	}

	server := service.NewServer(appEnv)
	server.Run(":" + port)

}
