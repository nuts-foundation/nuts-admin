package discovery

import (
	"context"
	"github.com/nuts-foundation/go-nuts-client/nuts"
	"github.com/nuts-foundation/go-nuts-client/nuts/discovery"
)

type Service struct {
	Client *discovery.Client
}

func (i Service) GetDiscoveryServices(ctx context.Context) ([]discovery.ServiceDefinition, error) {
	httpResponse, err := i.Client.GetServices(ctx)
	response, err := nuts.ParseResponse(err, httpResponse, discovery.ParseGetServicesResponse)
	if err != nil {
		return nil, err
	}
	return *response.JSON200, nil
}

func (i Service) ActivationStatus(ctx context.Context, serviceID string, subjectID string) (*DIDStatus, error) {
	httpResponse, err := i.Client.GetServiceActivation(ctx, serviceID, subjectID)
	response, err := nuts.ParseResponse(err, httpResponse, discovery.ParseGetServiceActivationResponse)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	result := &DIDStatus{
		ServiceID: serviceID,
		Active:    response.JSON200.Activated,
	}
	if response.JSON200.Vp != nil {
		result.Presentations = *response.JSON200.Vp
	}
	return result, nil
}
