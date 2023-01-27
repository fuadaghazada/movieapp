package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Registry defines a service registry.
type Registry interface {
	Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error
	Deregister(ctx context.Context, instanceID string, serviceName string) error
	ServiceAddresses(ctx context.Context, serviceID string) ([]string, error)
	ReportHealthyState(instanceID string, serviceName string) error
}

// ErrNotFound is returned when no service addresses are found
var ErrNotFound = errors.New("no service address found")

// GenerateInstanceID generates a pseudo-random  service instance identifier
func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
