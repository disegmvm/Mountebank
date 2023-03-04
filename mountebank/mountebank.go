package mountebank

import (
	"context"
	"github.com/senseyeio/mbgo"
	"net/http"
	"net/url"
)

type Mountebank struct {
	client *mbgo.Client
}

func NewMountebank() (*Mountebank, error) {
	httpClient := &http.Client{}
	mbURL, err := url.Parse("http://localhost:2525/")
	if err != nil {
		return &Mountebank{}, err
	}
	client := mbgo.NewClient(httpClient, mbURL)

	return &Mountebank{client: client}, nil
}

func (m Mountebank) CreateImposter(ctx context.Context, imposter mbgo.Imposter) error {
	_, err := m.client.Imposter(ctx, imposter.Port, false)
	if err != nil {
		_, err = m.client.Create(ctx, imposter)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m Mountebank) AddStub(ctx context.Context, port int, newStub mbgo.Stub) error {
	_, err := m.client.Imposter(ctx, port, true)
	if err != nil {
		return err
	}
	_, err = m.client.AddStub(ctx, port, -1, newStub)

	return err
}

func (m Mountebank) GetRequestsFromMountebank(ctx context.Context, port int) ([]mbgo.HTTPRequest, error) {
	res, err := m.client.Imposter(ctx, port, false)
	if err != nil {
		return nil, err
	}
	var requests []mbgo.HTTPRequest
	for _, request := range res.Requests {
		r := request.(*mbgo.HTTPRequest)
		requests = append(requests, *r)
	}

	return requests, nil
}
