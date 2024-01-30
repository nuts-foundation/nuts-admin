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

func (i Service) ActivationStatus(ctx context.Context, serviceID string, subjectID string) (*DIDStatus, error) {
	httpResponse, err := i.Client.GetServiceActivation(ctx, serviceID, subjectID)
	if err != nil {
		return nil, nuts.UnwrapAPIError(err)
	}
	response, err := ParseGetServiceActivationResponse(httpResponse)
	if err != nil {
		return nil, err
	}
	if response.JSON200 == nil {
		return nil, errors.New("unable to get service activation")
	}
	return &DIDStatus{
		ServiceID:    serviceID,
		Active:       response.JSON200.Activated,
		Presentation: response.JSON200.Vp,
	}, nil
}
