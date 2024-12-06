package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

func ReadFromRequestBody(request *http.Request, result interface{}) {
	log.Printf("%v", request.Body)
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	PanicIfError(err, "Failed to read from request body")
}

func WriteToResponseBody(writer http.ResponseWriter, response interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	PanicIfError(err, "Failed to write to response body")
}
