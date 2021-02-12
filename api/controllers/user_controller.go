package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"stockwatch/api/auth"
	"stockwatch/api/responses"
	"stockwatch/models"
)

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	err = user.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	_, err = user.Save(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, http.StatusCreated, map[string]bool{
		"success": true,
	})
}

func (server *Server) LoginUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	password := user.Password
	userInfo, err := user.FindByEmail(server.DB, user.Email)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = models.VerifyPassword(userInfo.Password, password)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if user.Password != userInfo.Password {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	signedToken, err := auth.CreateToken(userInfo.ID)

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusCreated, map[string]interface{}{
		"token":    signedToken,
		"email":    userInfo.Email,
		"username": userInfo.UserName,
		"id":       userInfo.ID,
	})
}
