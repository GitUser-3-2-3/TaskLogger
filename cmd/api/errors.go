package main

import (
	"fmt"
	"net/http"
)

func (bknd *backend) errResponseJSON(w http.ResponseWriter, r *http.Request, status int, errMsg any) {
	wrp := wrapper{"error": errMsg}
	err := bknd.writeJSON(w, status, wrp, nil)
	if err != nil {
		bknd.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (bknd *backend) errRateLimitExceeded(w http.ResponseWriter, r *http.Request) {
	errMsg := "rate limit exceeded, try again after a few seconds"
	bknd.errResponseJSON(w, r, http.StatusTooManyRequests, errMsg)
}

func (bknd *backend) errFailedValidation(w http.ResponseWriter, r *http.Request, errs map[string]string) {
	bknd.errResponseJSON(w, r, http.StatusBadRequest, errs)
}

func (bknd *backend) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	bknd.logger.Err(err).Msgf("method: %s, uri: %s", method, uri)
}

func (bknd *backend) errMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	errMsg := fmt.Sprintf("%s method not supported for this request", r.Method)
	bknd.errResponseJSON(w, r, http.StatusMethodNotAllowed, errMsg)
}

func (bknd *backend) errBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	errMsg := "invalid request"
	bknd.errResponseJSON(w, r, http.StatusBadRequest, errMsg+", "+err.Error())
}

func (bknd *backend) errResourceNotFound(w http.ResponseWriter, r *http.Request) {
	errMsg := "requested resource not found ·_·"
	bknd.errResponseJSON(w, r, http.StatusMethodNotAllowed, errMsg)
}

func (bknd *backend) errDuplicateEntryFound(w http.ResponseWriter, r *http.Request, err error) {
	errMsg := "duplicate entry not allowed"
	bknd.errResponseJSON(w, r, http.StatusMethodNotAllowed, errMsg+", "+err.Error())
}

func (bknd *backend) errInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	bknd.logError(r, err)
	errMsg := "server encountered a problem and could not process your request :("
	bknd.errResponseJSON(w, r, http.StatusInternalServerError, errMsg)
}
