package main

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"strings"
	"encoding/json"
	"bytes"
)

func TestInitPermutationHandler(t *testing.T) {
	userMap = InitMemMap()

	handler := http.HandlerFunc(initPermutation)


	var (
		validMethod = "POST"
		invalidMethod = "PATCH"

		validUserID= "1.1.1.1:1234"

		validInput = []int{1,2}
		emptyInput = "[]"
		invalidBody = "randomString"
	)

	// case 1: incorrect method
	req, err := http.NewRequest(invalidMethod, "/v1/init", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = validUserID
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("[case 1] handler returned wrong status code: expected %v got %v",
			http.StatusMethodNotAllowed, status)
	}

	// case 2: empty input
	req, err = http.NewRequest(validMethod, "/v1/init", strings.NewReader(emptyInput))
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = validUserID
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("[case 2] handler returned wrong status code: expected %v got %v",
			http.StatusBadRequest, status)
	}

	// case 3: invalid body content
	req, err = http.NewRequest(validMethod, "/v1/init", strings.NewReader(invalidBody))
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = validUserID
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("[case 3] handler returned wrong status code: expected %v got %v",
			http.StatusBadRequest, status)
	}

	// case 4: valid input
	inputBytes, err := json.Marshal(validInput)
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest(validMethod, "/v1/init", bytes.NewReader(inputBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = validUserID
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("[case 4] handler returned wrong status code: expected %v got %v",
			http.StatusOK, status)
	}


}

func TestGetNextPermutationHandler(t *testing.T) {
	userMap = InitMemMap()

	handler := http.HandlerFunc(getNextPermutation)

	var(
		validUserID = "0.0.0.0:1234"
		invalidUserID = "1.1.1.1:1234"
		samplePermutation = []int{1,2}

		validMethod = "GET"
		invalidMethod = "POST"
	)

	// init sample user
	userMap.InitUser(validUserID, samplePermutation)

	// case 1: incorrect method
	req, err := http.NewRequest(invalidMethod, "/v1/next", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("[case 1] handler returned wrong status code: expected %v got %v",
			http.StatusMethodNotAllowed, status)
	}

	// case 2: incorrect user
	req, err = http.NewRequest(validMethod, "/v1/next", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = invalidUserID
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("[case 2] handler returned wrong status code: expected %v got %v",
			http.StatusNotFound, status)
	}

	// case 3: valid user
	req, err = http.NewRequest(validMethod, "/v1/next", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.RemoteAddr = validUserID
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("[case 3] handler returned wrong status code: expected %v got %v",
			http.StatusOK, status)
	}

	// Check the response body is what we expect.
	expected, err := json.Marshal(samplePermutation)
	if err != nil {
		t.Fatal(err)
	}
	if rr.Body.String() != string(expected) {
		t.Errorf("handler returned unexpected body: expected %v, got %v",
			expected, rr.Body.String())
	}
}
