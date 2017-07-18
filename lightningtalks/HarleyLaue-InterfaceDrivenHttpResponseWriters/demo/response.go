package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// Responser is an interface that modifies the http.ResponseWriter and potentially
// transforms the data as the second parameter.
//
// This example code uses an interface instead of functions for intaking the
// ResponseWriter & data. Thus this doesn't strictly follow the "Functional options"
// pattern.
type Responser interface {
	Response(http.ResponseWriter, *interface{}) error
}

type (
	// Used to specify a content type & a marshaler for that type
	Accept struct {
		// This probably could have been of ContentType
		ContentType string
		// For the most part, this should conform to what json, xml, and probably
		// other libraries already use for marshalling data
		Marshaler func(interface{}) ([]byte, error)
	}
	// Used to specify a content-type in a response
	ContentType string
	// Error implements the error interface to allow it to be treated like an error
	// and to take over the response. This is a terminating Responser that will write
	// headers & content to the body
	Error struct {
		Error      error
		StatusCode int
		Accept     Accept
	}
	// Header is just a wrapper around http.Header to illistrate things that can be
	// done with a Responser
	Header http.Header
	// StatusCode is another instance of an example of modifying the ResponseWriter.
	// This implementation will actually defer writing the status code until it verifies
	// the other Responsers have not errored since the Error type writes its own
	// StatusCode
	StatusCode int
)

func (a Accept) Response(w http.ResponseWriter, v *interface{}) error {
	if a.ContentType != "" {
		if err := ContentType(a.ContentType).Response(w, nil); err != nil {
			return errors.WithStack(err)
		}
	}

	if a.Marshaler != nil {
		b, err := a.Marshaler(*v)
		if err != nil {
			return errors.WithStack(err)
		}

		*v = string(b)
	}

	return nil
}

func (c ContentType) Response(w http.ResponseWriter, _ *interface{}) error {
	w.Header().Set("Content-Type", string(c))
	return nil
}

func (e Error) Response(w http.ResponseWriter, v *interface{}) error {
	if e.Error == nil {
		return nil
	}

	w.WriteHeader(e.StatusCode)
	if e.Accept.Marshaler != nil {
		nv := interface{}(e.Error)
		*v = nv
		return e.Accept.Response(w, &nv)
	}

	return errors.WithStack(e.Error)
}

func (h Header) Response(w http.ResponseWriter, _ *interface{}) error {
	for k, v := range h {
		w.Header().Del(k)
		for _, v := range v {
			w.Header().Add(k, v)
		}
	}

	return nil
}

func (s StatusCode) Response(w http.ResponseWriter, _ *interface{}) error {
	w.WriteHeader(int(s))
	return nil
}

func (s StatusCode) StatusCode() int {
	return int(s)
}

// Acceptable looks at the request headers for the Accept header and will select the
// accept[n] that matches that content-type. If none match, it will return a 415
// and plaintext message
func Acceptable(r *http.Request, accepts ...Accept) Responser {
	accept := r.Header.Get("Accept")
	if (accept == "" || accept == "*/*") && len(accepts) > 0 {
		return accepts[0]
	}

	for _, a := range accepts {
		if strings.Contains(accept, a.ContentType) {
			return a
		}
	}

	err := Error{
		Error:      errors.New("unable to find an acceptable media type"),
		StatusCode: http.StatusUnsupportedMediaType,
	}
	if len(accepts) > 0 {
		err.Accept = accepts[0]
	}

	return err
}

// Write is an example implementation for taking interfaces as an argument for
// setting options on the ResponseWriter or transforming the data v. Responsers are
// done in order. This means the order you pass them to this function matters. For
// instance, a marshaller that converts a struct to a string to be output should be
// called last. An interface used to protect data on structs would likely be called
// first. This implementation also takes StatusCode's to defer writing those until
// the end in case the Error type is returned (in the case of Acceptable.)
func Write(w http.ResponseWriter, v interface{}, responses ...Responser) error {
	type statusCoder interface {
		StatusCode() int
	}

	var statusCode StatusCode
	for _, r := range responses {
		switch t := r.(type) {
		case StatusCode:
			statusCode = t
		case statusCoder:
			statusCode = StatusCode(t.StatusCode())
		default:
			if err := r.Response(w, &v); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	// defer writing the status code until after everything else since an
	// Error.Response, for example, may write a status code
	if statusCode != (StatusCode(0)) {
		statusCode.Response(w, nil)
	}

	_, err := fmt.Fprint(w, v)
	return errors.WithStack(err)
}
