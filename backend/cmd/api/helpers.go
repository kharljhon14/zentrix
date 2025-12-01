package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"strings"
)

type envelope map[string]any

func (app application) writeJSON(
	w http.ResponseWriter,
	status int,
	data envelope,
	headers http.Header) error {

	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	maps.Insert(w.Header(), maps.All(headers))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	w.Write(js)

	return nil
}

func (app application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formed JSON (at character %d)", syntaxError.Offset)

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at characted %d)", unmarshalTypeError.Offset)

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formed JSON")

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger that %d bytes", maxBytesError.Limit)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(dst)
	if !errors.Is(err, io.EOF) {
		return errors.New("body must contain a single JSON value")
	}

	return nil
}
