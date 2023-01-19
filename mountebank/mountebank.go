package mountebank

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/onsi/gomega"

	"github.com/onsi/gomega/matchers"

	"github.com/senseyeio/mbgo"
)

const (
	AllPaths = "all-paths"
	Port     = 8080
)

type Mountebank struct {
	client *mbgo.Client
}

func NewMountebank(env string) (*Mountebank, error) {
	httpClient := &http.Client{}
	mbURL, err := url.Parse(fmt.Sprintf("http://mountebank.%s.warehouse.ri-tech.io:2525/", env))
	if err != nil {
		return &Mountebank{}, err
	}
	client := mbgo.NewClient(httpClient, mbURL)

	return &Mountebank{client: client}, nil
}

func (m Mountebank) CreateImposters(ctx context.Context, imposters []mbgo.Imposter) error {
	//_, err := m.client.DeleteAll(ctx, true)
	//if err != nil {
	//	return err
	//}

	for _, imposter := range imposters {
		//FIXME Add error validation (should create only if imposter not found)
		_, err := m.client.Imposter(ctx, imposter.Port, false)
		if err != nil {
			_, err = m.client.Create(ctx, imposter)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m Mountebank) AddStub(ctx context.Context, port int, newStub mbgo.Stub) error {
	imposter, err := m.client.Imposter(ctx, port, true)
	if err != nil {
		return err
	}

	for i, stub := range imposter.Stubs {
		if reflect.DeepEqual(stub.Predicates, newStub.Predicates) {
			_, err = m.client.OverwriteStub(ctx, port, i, newStub)
			if err != nil {
				return err
			}

			return nil
		}
	}

	_, err = m.client.AddStub(ctx, port, -1, newStub)

	return err
}

func (m Mountebank) RemoveStub(ctx context.Context, port int, stubToRemove mbgo.Stub) error {
	imposter, err := m.client.Imposter(ctx, port, true)
	if err != nil {
		return err
	}

	for i, stub := range imposter.Stubs {
		if reflect.DeepEqual(stub.Predicates, stubToRemove.Predicates) {
			_, err = m.client.RemoveStub(ctx, port, i)

			return err
		}
	}

	return nil
}

func (m Mountebank) GetRecordedRequestsByPath(ctx context.Context, port int, path string) ([]mbgo.HTTPRequest, error) {
	res, err := m.client.Imposter(ctx, port, false)
	if err != nil {
		return nil, err
	}

	var requests []mbgo.HTTPRequest
	for _, request := range res.Requests {
		r := request.(*mbgo.HTTPRequest)
		if path == AllPaths || path == r.Path {
			requests = append(requests, *r)
		}
	}

	return requests, nil
}

func (m Mountebank) GetRecordedRequestsByStub(ctx context.Context, port int, stub mbgo.Stub) ([]mbgo.HTTPRequest, error) {
	return m.GetRecordedRequestsByPath(ctx, port, stub.Predicates[0].Request.(*mbgo.HTTPRequest).Path)
}

func (m Mountebank) GetRecordedRequestByHeaderValue(ctx context.Context, header string, value string, pollIntervalInSeconds int, timeoutInSeconds int) mbgo.HTTPRequest {
	reqs := m.GetRecordedRequestsByHeaderValue(
		ctx, header, value, 1, pollIntervalInSeconds, timeoutInSeconds)

	return reqs[0]
}

func (m Mountebank) GetRecordedRequestsByHeaderValue(ctx context.Context, header string, value string, expectedNumOfRequests int, pollIntervalInSeconds int, timeoutInSeconds int) []mbgo.HTTPRequest {
	//endTime := time.Now().UTC().Add(time.Second * time.Duration(timeoutInSeconds))
	var requests []mbgo.HTTPRequest
	found := make(map[string]struct{})

	gomega.Eventually(func() int {
		reqs, _ := m.GetRecordedRequestsByPath(ctx, Port, AllPaths)

		for _, req := range reqs {
			if req.Headers.Get(header) == value {
				//FIXME This line of code makes me sad -_-
				if _, ok := found[req.Body.(string)]; !ok {
					requests = append(requests, req)
					found[req.Body.(string)] = struct{}{}
				}
			}
		}

		return len(requests)
	}).
		WithPolling(time.Second*time.Duration(pollIntervalInSeconds)).
		WithTimeout(time.Second*time.Duration(timeoutInSeconds)).
		Should(gomega.Equal(expectedNumOfRequests),
			fmt.Sprintf("expected requests num: %d, found: %d; \n%v",
				expectedNumOfRequests, len(requests), reflect.ValueOf(found).MapKeys()))

	return requests
}

func (m Mountebank) GetRecordedRequestByBody(
	ctx context.Context,
	jsonBody string,
	pollIntervalInSeconds int,
	timeoutInSeconds int) (mbgo.HTTPRequest, error) {
	endTime := time.Now().UTC().Add(time.Second * time.Duration(timeoutInSeconds))
	matcher := matchers.MatchJSONMatcher{JSONToMatch: jsonBody}

	for time.Now().UTC().Before(endTime) {
		requests, _ := m.GetRecordedRequestsByPath(ctx, Port, AllPaths)

		for _, req := range requests {
			matches, _ := matcher.Match(req.Body.(string))
			if matches {
				return req, nil
			}
		}

		time.Sleep(time.Second * time.Duration(pollIntervalInSeconds))
	}

	return mbgo.HTTPRequest{}, fmt.Errorf("request not found")
}
