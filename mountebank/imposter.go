package mountebank

import "github.com/senseyeio/mbgo"

var Imposter = mbgo.Imposter{
	Port:           8181,
	Proto:          "http",
	Name:           "imposter",
	RecordRequests: true,
	AllowCORS:      true,
	DefaultResponse: mbgo.HTTPResponse{
		StatusCode: 200,
	},
	Stubs: []mbgo.Stub{},
}
