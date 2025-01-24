package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JSONResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

func (app *application) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, val := range headers[0] {
			w.Header()[key] = val
		}
	}

	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(status)
	_, err = w.Write(payload)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1 << 20
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// make sure there is only one json file in the request body
	// using a throwaway variable like empty struct
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body should contain only one json value")
	}
	return nil
}

func (app *application) ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var jsonResponse JSONResponse
	jsonResponse.Error = true
	jsonResponse.Message = err.Error()

	return app.WriteJSON(w, statusCode, jsonResponse)
}