package response

import (
	"github.com/go-chi/render"
	"net/http"
)

type Response struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

func Send200Success(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
}

func Send400Error(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, Response{Error: "very bad"})
}

func Send401Error(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusUnauthorized)
	render.JSON(w, r, Response{Error: "invalid token("})
}

func Send404Error(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, Response{Error: "resource not found"})
}

func Send403Error(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusForbidden)
	render.JSON(w, r, Response{Error: "no access"})
}

func Send500Error(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, Response{Error: "technical chocolates..."})
}
