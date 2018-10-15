package main

// The point of program is to create new permutation from int array for POST request to defined endpoint (/v1/init)
// and via GET request to another endpoint (/v1/next) receive next permutation for initial array.
// Separate users will be defined via IP (included port)
// User's initial array will be redefined with new POST to init endpoint


// CAUTION! Throughout this example unsupported by json.Marshal values (cyclic, channels, infinity values etc.) are NOT USED,
// so error checks for this methods can be skipped for code clarity and simplicity.
// However it is strongly not recommended to do this in production code.

import (
	"net/http"

	"github.com/gorilla/mux"
	"io/ioutil"
	"encoding/json"
)

type UserRecords  interface {
	InitUser(userID string, permutation []int)
	UserExists(userID string) bool
	NextForUser(userID string) (permutation []int, userExists bool)
}

var userMap UserRecords

func main() {
	userMap = InitMemMap()

	r := mux.NewRouter()
	r.HandleFunc("/v1/init", initPermutation).Methods("POST")
	r.HandleFunc("/v1/next", getNextPermutation).Methods("GET")

	http.ListenAndServe(":8080", r)
}

// initPermutation initialize permutation for user
func initPermutation(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errMsg, _ := json.Marshal(NewErrorMsg(http.StatusInternalServerError, ErrorDataReading))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errMsg)
		return
	}
	defer r.Body.Close()

	var input []int
	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		errMsg, _ := json.Marshal(NewErrorMsg(http.StatusBadRequest, ErrorInvalidInput))
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errMsg)
		return
	}

	if len(input) == 0 {
		errMsg, _ := json.Marshal(NewErrorMsg(http.StatusBadRequest, ErrorInvalidInput))
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errMsg)
		return
	}

	userID := r.RemoteAddr
	userMap.InitUser(userID, input)
	w.WriteHeader(http.StatusOK)
}

// getNextPermutation retrieves next permutation for current user
func getNextPermutation(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID := r.RemoteAddr
	userExists := userMap.UserExists(userID)
	if !userExists {
		errMsg, _ := json.Marshal(NewErrorMsg(http.StatusNotFound, ErrorPermNotInitialized))
		w.WriteHeader(http.StatusNotFound)
		w.Write(errMsg)
		return
	}

	perm, _ := userMap.NextForUser(userID)

	bytes, _ := json.Marshal(perm)

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}