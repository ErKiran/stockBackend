package api

import (
	"fmt"
	"os"
	"stockwatch/api/controllers"
)

var server = controllers.Server{}

func Run() {
	port := fmt.Sprintf(":%v", os.Getenv("HTTP_SERVER_PORT"))
	server.Run(port)
}
