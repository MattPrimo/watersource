package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Sample struct {
	Id    int    `json:"id"`
	What  string `json:"what"`
	Who   string `json:"who"`
	When  string `json:"when"`
	Notes string `json:"notes"`
}

var idCounter int
var samples = []Sample{
	{
		What:  "Cromwell Lake",
		Who:   "Outlet Stream",
		When:  "October 2013",
		Notes: "Protected lake; within Montreal University field station",
	},
}

func SamplesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(samples)
	if err != nil {
		panic(err)
	}
	w.Write(j)
}

func CreateSampleHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming Sample from the request body
	var sample Sample
	err := json.NewDecoder(r.Body).Decode(&sample)
	if err != nil {
		panic(err)
	}

	idCounter++

	// Grab the sample and set some dummy data
	sample.Id = idCounter

	samples = append(samples, sample)

	// Serialize the modified sample to JSON
	j, err := json.Marshal(sample)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func UpdateSampleHandler(w http.ResponseWriter, r *http.Request) {
	// Grab the sample's id from the incoming url
	var sample Sample
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}

	// Decode the incoming sample json
	err = json.NewDecoder(r.Body).Decode(&sample)
	if err != nil {
		panic(err)
	}

	// Find the sample in our samples slice and upate it's name
	for index, _ := range samples {
		if samples[index].Id == id {
			samples[index] = sample
		}
	}

	// Respond with a 204 indicating success, but no content
	w.WriteHeader(http.StatusNoContent)
}

func DeleteSampleHandler(w http.ResponseWriter, r *http.Request) {
	// Grab the sample's id from the incoming url
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}

	// Find the index of the sample
	sampleIndex := -1
	for index, _ := range samples {
		if samples[index].Id == id {
			sampleIndex = index
			break
		}
	}

	// If we actually found a sample remove it from the slice
	if sampleIndex != -1 {
		samples = append(samples[:sampleIndex], samples[sampleIndex+1:]...)
	}

	// Respond with a 204 indicating success, but no content
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	log.Println("Starting Server")

	r := mux.NewRouter()
	r.HandleFunc("/api/samples", SamplesHandler).Methods("GET")
	r.HandleFunc("/api/samples", CreateSampleHandler).Methods("POST")
	r.HandleFunc("/api/samples/{id}", UpdateSampleHandler).Methods("PUT")
	r.HandleFunc("/api/samples/{id}", DeleteSampleHandler).Methods("DELETE")
	http.Handle("/api/", r)

	http.Handle("/", http.FileServer(http.Dir("./")))

	log.Println("Listening on 3030")
	http.ListenAndServe(":3030", nil)
}
