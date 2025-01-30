package v2

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

func validateResponse(requests []ResponseRequests) error {
	var result error
	for _, item := range requests {
		if item.State == "INVALID" {
			result = multierror.Append(result, fmt.Errorf("%+v", item))
		}
	}

	return result
}

var client = http.Client{}

type RespErrors struct {
	Errors []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

func jsonReq[T any](req *http.Request) (*T, error) {
	response, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "http.DefaultClient.Do")
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("invalid credentials")
	}

	var s T
	payload, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "io.ReadAll")
	}

	if response.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("bad request: %s", string(payload))
	}

	var respErr RespErrors
	if err := json.Unmarshal(payload, &respErr); err == nil && len(respErr.Errors) > 0 {
		return nil, fmt.Errorf("json error: %v", respErr)
	}

	if err := json.Unmarshal(payload, &s); err != nil {
		return nil, fmt.Errorf("json error: %s", payload)
	}

	return &s, nil
}
