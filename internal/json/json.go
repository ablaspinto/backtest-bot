package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	MaxBodySize = 1 << 20 // 1 MB

	contentTypeJSON = "application/json"
)

var (
	ErrBodyTooLarge     = errors.New("request body too large")
	ErrEmptyBody        = errors.New("request body is empty")
	ErrMalformedJSON    = errors.New("malformed JSON")
	ErrUnknownField     = errors.New("unknown field in JSON")
	ErrMultipleJSONObjs = errors.New("request body must contain only one JSON object")
)

func Write(w http.ResponseWriter, status int, data any) error {
	return WriteWithHeaders(w, status, data, nil)
}

func WriteWithHeaders(w http.ResponseWriter, status int, data any, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("json encode error: %w", err)
	}

	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(status)

	_, err = w.Write(js)
	return err
}

func Read(r *http.Request, data any) error {
	return ReadWithMaxSize(r, data, MaxBodySize)
}
func ReadOrError[T any](w http.ResponseWriter, r *http.Request) (T, bool) {
	var data T
	if err := Read(r, &data); err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return data, false
	}
	return data, true
}

func ReadWithMaxSize(r *http.Request, data any, maxSize int64) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != "" && !strings.HasPrefix(contentType, contentTypeJSON) {
		return fmt.Errorf("content-type must be %s", contentTypeJSON)
	}

	r.Body = http.MaxBytesReader(nil, r.Body, maxSize)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(data); err != nil {
		return categorizeDecodeError(err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return ErrMultipleJSONObjs
	}

	return nil
}

func categorizeDecodeError(err error) error {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var maxBytesError *http.MaxBytesError

	switch {
	case errors.As(err, &syntaxError):
		return fmt.Errorf("%w at position %d", ErrMalformedJSON, syntaxError.Offset)

	case errors.As(err, &unmarshalTypeError):
		return fmt.Errorf("invalid type for field %q: expected %s",
			unmarshalTypeError.Field, unmarshalTypeError.Type)

	case errors.As(err, &maxBytesError):
		return ErrBodyTooLarge

	case errors.Is(err, io.EOF):
		return ErrEmptyBody

	case strings.HasPrefix(err.Error(), "json: unknown field"):
		return fmt.Errorf("%w: %s", ErrUnknownField, strings.TrimPrefix(err.Error(), "json: unknown field "))

	default:
		return fmt.Errorf("%w: %v", ErrMalformedJSON, err)
	}
}

func Error(w http.ResponseWriter, status int, message string) error {
	return Write(w, status, map[string]string{"error": message})
}

func WriteSuccess(w http.ResponseWriter, status int, data any) error {
	return Write(w, status, Envelope{
		Success: true,
		Data:    data,
	})
}
func WriteError(w http.ResponseWriter, status int, message string) error {
	return Write(w, status, Envelope{
		Success: false,
		Error:   message,
	})
}
