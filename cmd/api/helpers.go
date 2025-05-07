package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type wrapper map[string]any

func (bknd *backend) writeJSON(w http.ResponseWriter, status int, data wrapper, headers http.Header) error {
	response, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	response = append(response, byte('\n'))
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(response)
	if err != nil {
		bknd.logger.Err(err).Msg("Error while writing response")
		return err
	}
	return nil
}

func (bknd *backend) readIdParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}
