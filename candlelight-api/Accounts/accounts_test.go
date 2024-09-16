package Accounts

import (
	"bytes"
	"candlelight-ruleengine/Accounts"
	"candlelight-ruleengine/Engine"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_CreateAccount(t *testing.T) {

	tests := []struct {
		name               string
		expectedStatusCode int
		submission         interface{}
	}{
		{
			name:               "Valid New Account",
			expectedStatusCode: http.StatusOK,
			submission: Accounts.User{
				Username: Engine.GenerateId(), //Using generateId to get a guaranteed unique username
				Password: "password",
			},
		},
		{
			name:               "Pre-Existing Account",
			expectedStatusCode: http.StatusBadRequest,
			submission: Accounts.User{
				Username: "ryan",
				Password: "pass123",
			},
		},
		{
			name:               "Malformed User Object",
			expectedStatusCode: http.StatusBadRequest,
			submission: struct {
				field  string
				field2 string
			}{
				field:  "asdf",
				field2: "asdf",
			},
		},
		{
			name:               "Missing User Object",
			expectedStatusCode: http.StatusBadRequest,
			submission:         nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.submission)

			request := httptest.NewRequest(http.MethodPost, "/createAccount", bytes.NewReader(body))
			response := httptest.NewRecorder()

			CreateAccount(response, request)

			if response.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Status Code Error! Expected {%d} but received {%d}", tt.expectedStatusCode, response.Result().StatusCode)
			}

			returnedUser := Accounts.SafeUser{}

			err := json.Unmarshal(response.Body.Bytes(), &returnedUser)
			if tt.expectedStatusCode == http.StatusOK {
				if err != nil {
					t.Errorf("Error with unmarshalling User. Should not have error, but got {%s}", err)
				}

				sentUser := tt.submission.(Accounts.User)

				if returnedUser.Username != sentUser.Username {
					t.Errorf("Error with returned User. Expected username {%s} but got username {%s}", sentUser.Username, returnedUser.Username)
				}

				//Clean up created user object
				Engine.RDB.Del(Engine.RDB.Context(), "user:"+sentUser.Username)
			} else {
				if err == nil {
					t.Error("Error unmarshalling User. Should have error, but got none!")
				}
			}
		})
	}
}

func Test_Login(t *testing.T) {
	tests := []struct {
		name               string
		expectedStatusCode int
		submission         interface{}
	}{
		{
			name:               "Valid Login", //Assumes a user with Username "ryan" and Password "pass123" already exists
			expectedStatusCode: http.StatusOK,
			submission: Accounts.User{
				Username: "ryan",
				Password: "pass123",
			},
		},
		{
			name:               "Wrong Password",
			expectedStatusCode: http.StatusUnauthorized,
			submission: Accounts.User{
				Username: "ryan",
				Password: "wrongpassword",
			},
		},
		{
			name:               "Non-Existant User",
			expectedStatusCode: http.StatusUnauthorized,
			submission: Accounts.User{
				Username: Engine.GenerateId(), //Using GenerateId to hopefully get a non-existant username
				Password: "password",
			},
		},
		{
			name:               "Malformed User Object",
			expectedStatusCode: http.StatusBadRequest,
			submission: struct {
				field1 string
				field2 string
			}{
				field1: "asdf",
				field2: "asdf",
			},
		},
		{
			name:               "Missing User Object",
			expectedStatusCode: http.StatusBadRequest,
			submission:         nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.submission)

			request := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			response := httptest.NewRecorder()

			Login(response, request)

			if response.Result().StatusCode != tt.expectedStatusCode {
				t.Errorf("Status Code Error! Expected {%d} but received {%d}", tt.expectedStatusCode, response.Result().StatusCode)
			}

			returnedUser := Accounts.SafeUser{}

			err := json.Unmarshal(response.Body.Bytes(), &returnedUser)
			if tt.expectedStatusCode == http.StatusOK {
				if err != nil {
					t.Errorf("Error with unmarshalling User. Should not have error, but got {%s}", err)
				}

				sentUser := tt.submission.(Accounts.User)

				if returnedUser.Username != sentUser.Username {
					t.Errorf("Error with returned User. Expected username {%s} but got username {%s}", sentUser.Username, returnedUser.Username)
				}
			} else {
				if err == nil {
					t.Error("Error unmarshalling User. Should have error, but got none!")
				}
			}
		})
	}
}
