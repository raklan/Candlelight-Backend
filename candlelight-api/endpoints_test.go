package main

import (
	"bytes"
	"candlelight-api/CreationStudio"
	"candlelight-models/Game"
	"candlelight-ruleengine/Engine"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Studio_GET(t *testing.T) {
	ensureDummyGameExists()
	tests := []struct {
		name               string
		queryString        string
		expectedStatusCode int
	}{
		{
			name:               "Valid Game ID",
			queryString:        "?id=game123",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid Game ID",
			queryString:        "?id=invalid",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Missing Game ID",
			queryString:        "?id=",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Only Invalid Query Param",
			queryString:        "?invalid=game123",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "No Query Param",
			queryString:        "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Extra Invalid Query Param",
			queryString:        "?id=game123&invalid=asdf",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/studio"+tt.queryString, nil)
			response := httptest.NewRecorder()

			CreationStudio.Studio(response, request)

			//Check status code
			if response.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Status Code mismatch! Expected {%d} but got {%d}", tt.expectedStatusCode, response.Result().StatusCode)
			}

			//IFF success, body should be able to deserialize to a game object whose ID matches given game ID
			game := Game.Game{}
			err := json.Unmarshal(response.Body.Bytes(), &game)

			if tt.expectedStatusCode == http.StatusOK {
				if err != nil {
					t.Errorf("Error unmarshalling game! %s", err.Error())
				}
				if game.Id != "game123" {
					t.Errorf("Error with returned game! Expected a game with id {game123} but got {%s}", game.Id)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error during unmarshalling but none occurred!")
				}
				if game.Id == "game123" {
					t.Error("Error with returned game! Should not have gotten game with id {game123} but one was returned!")
				}
			}
		})
	}
}

func Test_Studio_POST(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatusCode int
		gameToSave         Game.Game
	}{
		{
			name:               "Valid Save Request",
			expectedStatusCode: 200,
			gameToSave: Game.Game{
				Name: "Valid Save Request_Test Game",
			},
		},
		{
			name:               "No Body",
			expectedStatusCode: 400,
			gameToSave:         Game.Game{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/studio", nil)
			if tt.name != "No Body" {
				asJson, err := json.Marshal(tt.gameToSave)
				if err != nil {
					t.Fatalf("Error trying to marshal tt.gameToSave: %s", err)
				}
				request = httptest.NewRequest(http.MethodPost, "/studio", bytes.NewReader(asJson))
			}
			response := httptest.NewRecorder()

			CreationStudio.Studio(response, request)

			if response.Result().StatusCode != tt.expectedStatusCode {
				t.Fatalf("Unexpected Status Code returned! Expected {%d}, Got {%d}", tt.expectedStatusCode, response.Result().StatusCode)
			}

			returned := Game.Game{}
			err := json.Unmarshal(response.Body.Bytes(), &returned)
			if tt.expectedStatusCode == 200 {
				if err != nil {
					t.Fatalf("Unexpected error trying to unmarshal response into Game object! %s", err)
				}

				if returned.Name != tt.gameToSave.Name {
					t.Errorf("Returned Game's name doesn't match given game! Expected: {%s}, Got {%s}", tt.gameToSave.Name, returned.Name)
				}
			} else {
				if err == nil {
					t.Fatalf("Expected error trying to unmarshal response into Game object, but got none!")
				}
			}

			if returned.Id != "" {
				Engine.RDB.Del(Engine.RDB.Context(), "game:"+returned.Id)
			}
		})
	}
}

func Test_Studio_DELETE(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatusCode int
		queryString        string
		idToDelete         string
	}{
		{
			name:               "Valid Game ID",
			queryString:        "?id=game123",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid Game ID",
			queryString:        "?id=invalid",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "Missing Game ID",
			queryString:        "?id=",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Only Invalid Query Param",
			queryString:        "?invalid=game123",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "No Query Param",
			queryString:        "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Extra Invalid Query Param",
			queryString:        "?id=game123&invalid=asdf",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ensureDummyGameExists()
			request := httptest.NewRequest(http.MethodDelete, "/studio"+tt.queryString, nil)
			response := httptest.NewRecorder()

			CreationStudio.Studio(response, request)

			//Check status code
			if response.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Status Code mismatch! Expected {%d} but got {%d}", tt.expectedStatusCode, response.Result().StatusCode)
			}

			//IFF success, body should be able to deserialize to a game object whose ID matches given game ID
			deleteResponse := struct {
				DeletedId string
			}{}
			err := json.Unmarshal(response.Body.Bytes(), &deleteResponse)

			if tt.expectedStatusCode == http.StatusOK {
				if err != nil {
					t.Errorf("Error unmarshalling response! %s", err.Error())
				}
				if deleteResponse.DeletedId != "game123" {
					t.Errorf("Error with returned ID! Expected id {game123} but got {%s}", deleteResponse.DeletedId)
				}
			} else {
				if err == nil {
					t.Errorf("Expected error during unmarshalling but none occurred!")
				}
				if deleteResponse.DeletedId == "game123" {
					t.Error("Error with returned game! Should not have gotten game with id {game123} but one was returned!")
				}
			}
		})
	}
}

// Generates the dummy game, if it doesn't already exist
func ensureDummyGameExists() {
	request := httptest.NewRequest(http.MethodGet, "/dummy", nil)
	response := httptest.NewRecorder()

	generateJSON(response, request)
}
