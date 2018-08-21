package views

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"regexp"

	"net/http"

	"webgo/models"
)

func init() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/movies1", AllMovies)
	http.HandleFunc("/movies2", CreateMovie)
	http.HandleFunc("/movies3", UpdateMovie)
	http.HandleFunc("/movies4", DeleteMovie)
	http.HandleFunc("/movies5/{id}", FindMovie)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World!")
}

func AllMovies(w http.ResponseWriter, r *http.Request) {
	movie := models.Movie{}
	movies, err := movie.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, movies)
}

func FindMovie(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile("/movies/(.*)")
	params := re.FindString(r.URL.Path)
	movie := models.Movie{}
	movie, err := movie.FindById(params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, movie)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	movie := models.Movie{}
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	movie.ID = bson.NewObjectId()
	if err := movie.Insert(movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, movie)
}

func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	movie := models.Movie{}
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	if err := movie.Update(movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, movie)
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	movie := models.Movie{}
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer r.Body.Close()

	if err := movie.Delete(movie); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, movie)
}
