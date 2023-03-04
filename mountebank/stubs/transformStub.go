package stubs

import (
	"github.com/senseyeio/mbgo"
)

var TransformStub = mbgo.Stub{
	Predicates: []mbgo.Predicate{
		{
			Operator: "equals",
			Request: &mbgo.HTTPRequest{
				Method: "POST",
				Path:   "/transform",
			},
		},
	},
	Responses: []mbgo.Response{
		{
			Type: "is",
			Value: mbgo.HTTPResponse{
				StatusCode: 200,
			},
		},
	},
}
