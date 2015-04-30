package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Sample struct {
	Id              int     `json:"id"`
	Title           string  `json:"title"`
	Location        string  `json:"location"`
	Date            string  `json:"date"`
	Notes           string  `json:"notes"`
	Temperature     int     `json:"temperature"`
	Width           float64 `json:"width"`
	Depth           float64 `json:"depth"`
	PH              float64 `json:"pH"`
	DissolvedOxygen float64 `json:"dissolvedOxygen"`
	Conductivity    float64 `json:"conductivity"`
	ORP             float64 `json:"ORP"`
	Discharge       float64 `json:"discharge"`
}

var idCounter = 1
var samples = []Sample{
	{
		Id:              0,
		Title:           "Bear Lake Site",
		Location:        "Bear Lake",
		Date:            "April 26th",
		Notes:           "Bear Lake was cool mate.",
		Temperature:     12,
		Width:           14.5,
		Depth:           14.5,
		PH:              14.5,
		DissolvedOxygen: 14.5,
		Conductivity:    14.5,
		ORP:             14.5,
		Discharge:       14.5,
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
