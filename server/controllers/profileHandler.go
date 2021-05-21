package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"paujim/auroraserverless/server/entities"
	"paujim/auroraserverless/server/repositories"
	"regexp"

	"github.com/ttacon/libphonenumber"
)

type InsertProfileResponse struct {
	ProfileID *int64 `json:"profile_id"`
}
type GetProfileResponse struct {
	Profiles []entities.Profile `json:"profiles"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func ProfileHandler(h repositories.SqlRepository) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}

		switch r.Method {
		case http.MethodGet:
			log.Printf("HTTP GET\n")
			// get data from DB
			profiles, err := h.GetProfiles()
			if err != nil {
				log.Printf("Error fetching data: %s\n", err)
				returnJSON(w, ErrorResponse{Error: "Error fetching data"}, http.StatusServiceUnavailable)
				return
			}

			// Return json
			log.Printf("Response: %v\n", profiles)
			returnJSON(w, GetProfileResponse{Profiles: profiles}, http.StatusOK)

		case http.MethodPost:
			log.Printf("HTTP POST\n")
			var profile entities.Profile
			err := json.NewDecoder(r.Body).Decode(&profile)
			// Validate
			if err != nil {
				log.Printf("Error decoding: %s\n", err)
				returnJSON(w, ErrorResponse{Error: "Bad Request"}, http.StatusBadRequest)
				return
			}
			log.Printf("Request Body: %v\n", profile)
			if !IsValidEmail(profile.Email) {
				log.Printf("Not valid email: %s\n", profile.Email)
				returnJSON(w, ErrorResponse{Error: fmt.Sprintf("Invalid email [%s].", profile.Email)}, http.StatusBadRequest)
				return
			}

			num, err := libphonenumber.Parse(profile.PhoneNumber, "AU")
			if err != nil || !libphonenumber.IsValidNumber(num) {
				log.Printf("Not a valid phone: %s\n", profile.PhoneNumber)
				returnJSON(w, ErrorResponse{Error: fmt.Sprintf("Invalid phone [%s].", profile.PhoneNumber)}, http.StatusBadRequest)
				return
			}
			formatedPhone := libphonenumber.Format(num, libphonenumber.E164)

			// Save to DB
			profileID, err := h.InsertProfile(profile.FullName, profile.Email, formatedPhone)
			if err != nil {
				log.Printf("Error inserting: %s\n", err)
				returnJSON(w, ErrorResponse{Error: "Unable to add the profile."}, http.StatusBadRequest)
				return
			}
			// Return json
			returnJSON(w, InsertProfileResponse{ProfileID: profileID}, http.StatusCreated)
		default:
			returnJSON(w, ErrorResponse{Error: "Only GET and POST methods are supported."}, http.StatusBadRequest)
		}
	}
}

func IsValidEmail(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(email)
}

func returnJSON(w http.ResponseWriter, jsonObj interface{}, statusCode int) {
	js, err := json.Marshal(jsonObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(js)
}
