package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"chatsmith-mini/pkg/failure"
)

const (
	headerAccept = "Accept"
	contentJSON  = "application/json"
)

var errProviderStatus = errors.New("oauth provider returned a non-200 status")

func getJSON(ctx context.Context, client *http.Client, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return failure.ErrWithTrace(err)
	}
	req.Header.Set(headerAccept, contentJSON)

	resp, err := client.Do(req)
	if err != nil {
		return failure.ErrWithTrace(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return failure.ErrWithTrace(errProviderStatus)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return failure.ErrWithTrace(err)
	}

	return nil
}
