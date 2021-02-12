package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"stockwatch/api/auth"
	"stockwatch/api/responses"
	"stockwatch/models"
	"strconv"

	"github.com/gorilla/mux"
)

func (server *Server) CreateContact(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	contact := models.Contact{}
	err = json.Unmarshal(body, &contact)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	contact.Prepare()
	err = contact.Validate()

	userID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	contact.UserID = int64(userID)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	savedContact, err := contact.Create(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"contact": savedContact,
	})
}

func (server *Server) UpdateContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	contactID, err := strconv.ParseInt(vars["id"], 10, 64)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	contact := models.Contact{}
	err = json.Unmarshal(body, &contact)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	_, err = auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	updatedContact, err := contact.Update(server.DB, contactID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"contact": updatedContact,
	})
}

func (server *Server) GetContactInfoOfUser(w http.ResponseWriter, r *http.Request) {
	contact := models.Contact{}

	userID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	contactInfo, err := contact.FindContactByUserID(server.DB, int64(userID))
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	responses.JSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"contact": contactInfo,
	})
}
