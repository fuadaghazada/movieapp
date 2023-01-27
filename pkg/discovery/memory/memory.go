package memory

import (
	"context"
	"movieexample.com/pkg/discovery"
	"sync"
	"time"
)

type Registry struct {
	sync.RWMutex
	serviceAddrs map[string]map[string]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

func NewRegistry() *Registry {
	return &Registry{
		serviceAddrs: map[string]map[string]*serviceInstance{},
	}
}

func (r *Registry) Register(_ context.Context, instID string, serviceName string, hostPort string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[serviceName]; !ok {
		r.serviceAddrs[serviceName] = map[string]*serviceInstance{}
	}
	r.serviceAddrs[serviceName][instID] = &serviceInstance{
		hostPort:   hostPort,
		lastActive: time.Now(),
	}

	return nil
}

func (r *Registry) Deregister(_ context.Context, instID string, svcName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.serviceAddrs[svcName]; !ok {
		return nil
	}
	delete(r.serviceAddrs[svcName], instID)
	return nil
}

func (r *Registry) ServiceAddresses(_ context.Context, svcName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.serviceAddrs[svcName]) == 0 {
		return nil, discovery.ErrNotFound
	}

	var res []string
	for _, sv := range r.serviceAddrs[svcName] {
		if sv.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, sv.hostPort)
	}

	return res, nil
}
