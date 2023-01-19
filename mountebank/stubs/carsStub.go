package stubs

import (
	"github.com/senseyeio/mbgo"
	"time"
)

type APIResponseHeader struct {
	DateTimeStamp time.Time `json:"dateTimeStamp"`
}

type APIResponseBody struct {
	Header  APIResponseHeader `json:"header"`
	Success bool              `json:"success"`
}

var CarsStub = mbgo.Stub{
	Predicates: []mbgo.Predicate{
		{
			Operator: "equals",
			Request: &mbgo.HTTPRequest{
				Method: "POST",
				Path:   "/test/cars",
			},
		},
	},
	Responses: []mbgo.Response{
		{
			Type: "is",
			Value: mbgo.HTTPResponse{
				StatusCode: 200,
				Body: APIResponseBody{
					Success: true,
				},
			},
		},
	},
}
