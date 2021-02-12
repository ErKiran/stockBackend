package controllers

import (
	"fmt"
	"log"
	"net/http"
	"stockwatch/models"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type Server struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (server *Server) Initialize() {
	server.Router = mux.NewRouter()
	server.DB = models.Init()
	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	server.Initialize()
	fmt.Println("Server is running on PORT :5000")
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
