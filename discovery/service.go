package discovery

import (
	"context"
	"errors"
	"github.com/nuts-foundation/nuts-admin/nuts"
)

type Service struct {
	Client *Client
}

func (i Service) GetDiscoveryServices(ctx context.Context) ([]ServiceDefinition, error) {
	httpResponse, err := i.Client.GetServices(ctx)
	if err != nil {
		return nil, nuts.UnwrapAPIError(err)
	}
	response, err := ParseGetServicesResponse(httpResponse)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, errors.New("unable to get services")
	}
	return *response.JSON200, nil
}
