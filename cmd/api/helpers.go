package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
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

func (bknd *backend) readIdParam(r *http.Request, param string) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName(param), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func (bknd *backend) readUUIDParam(r *http.Request, param string) (string, error) {
	params := httprouter.ParamsFromContext(r.Context())
	uid := params.ByName(param)
	if uid == "" || !isValidUUID(uid) {
		return "", errors.New("invalid uid")
	}
	return uid, nil
}

func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

func (bknd *backend) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dcr := json.NewDecoder(r.Body)
	dcr.DisallowUnknownFields()

	err := dcr.Decode(dst)
	if err != nil {
		return bknd.decodeJSONError(err)
	}
	err = dcr.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must contain exactly one JSON object")
	}
	return nil
}

func (bknd *backend) withTransaction(fn func(tx *sql.Tx) error) error {
	tx, err := bknd.models.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if nil == err {
			return
		}
		if val := recover(); val != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				bknd.logger.Error().Err(rollbackErr).Msg("Failed to rollback transaction after panic")
			}
			bknd.logger.Error().Interface("panic", val).Msg("Panic occurred during transaction")
			panic(val)
		} else if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				bknd.logger.Error().Err(err).Err(rollbackErr).Msg("Failed to rollback transaction")
			}
		}
	}()
	if err = fn(tx); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return tx.Commit()
}

func (bknd *backend) decodeJSONError(err error) error {
	var unmarshalTypeError *json.UnmarshalTypeError
	var syntaxError *json.SyntaxError
	var invalidUnmarshalError *json.InvalidUnmarshalError
	var maxBytesError *http.MaxBytesError

	switch {
	case errors.As(err, &syntaxError):
		return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
	case errors.Is(err, io.ErrUnexpectedEOF):
		return fmt.Errorf("body contains badly-formed JSON")
	case errors.Is(err, io.EOF):
		return fmt.Errorf("body must not be empty")
	case errors.As(err, &unmarshalTypeError):
		if unmarshalTypeError.Field != "" {
			return fmt.Errorf(
				"body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
		}
		return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)
	case strings.HasPrefix(err.Error(), "json: unknown field "):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return fmt.Errorf("body contains unknown JSON key %s", fieldName)
	case errors.As(err, &maxBytesError):
		return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)
	case errors.As(err, &invalidUnmarshalError):
		panic(err)
	default:
		return err
	}
}
