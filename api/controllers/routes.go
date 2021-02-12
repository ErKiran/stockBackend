package controllers

import (
	"net/http"
	"stockwatch/api/middleware"
)

func (server *Server) setJSON(path string, next func(http.ResponseWriter, *http.Request), method string) {
	server.Router.HandleFunc(path, middleware.SetMiddlewareJSON(next)).Methods(method, "OPTIONS")
}

func (server *Server) initializeRoutes() {
	server.Router.Use(middleware.CORS)
	server.Router.HandleFunc("/stock", server.CreateStockWatch).Methods("POST", "OPTIONS")
	server.Router.HandleFunc("/register", server.CreateUser).Methods("POST", "OPTIONS")
	server.Router.HandleFunc("/login", server.LoginUser).Methods("POST", "OPTIONS")
	server.Router.HandleFunc("/contact", server.CreateContact).Methods("POST", "OPTIONS")
	server.Router.HandleFunc("/contact/{id}", server.UpdateContact).Methods("PUT", "OPTIONS")
	server.Router.HandleFunc("/stock", server.GetStockScript).Methods("GET", "OPTIONS")
	server.Router.HandleFunc("/stocks", server.GetStockWatchOfUser).Methods("GET", "OPTIONS")
	server.Router.HandleFunc("/stock/{id}", server.GetStockDetail).Methods("GET", "OPTIONS")
	server.Router.HandleFunc("/contact", server.GetContactInfoOfUser).Methods("GET", "OPTIONS")
}
