package utils

import (
	"log"
	"sync"
)

type Service interface {
	Name() string
	Start()
	Stop()
	IsRunning() bool
}

// ServiceRegistry keeps track of all services
type ServiceRegistry struct {
	services map[string]Service
	mu       sync.Mutex
}

// NewServiceRegistry create a new Sevice Registry
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]Service, 10),
	}
}

// Register your service with the registry
func (r *ServiceRegistry) Register(s Service) Service {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.services[s.Name()] = s

	log.Printf("Registered service: '%v'", s.Name())

	return s
}

// Deregister your service from the registry
func (r *ServiceRegistry) Deregister(s Service) Service {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.services, s.Name())

	log.Printf("Deregistered service: '%v'", s.Name())

	return s
}

// DeregisterAll all your services from the registry
func (r *ServiceRegistry) DeregisterAll() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.services = make(map[string]Service, 10) // just trash the old map (best practice!)
}

// StartAll start all service
func (r *ServiceRegistry) StartAll() {
	var wg sync.WaitGroup
	wg.Add(len(r.services))

	for _, s := range r.services {
		go func() {
			defer wg.Done()

			log.Printf("Starting service: '%v'...", s.Name())

			s.Start()

			log.Printf("Started service: '%v'.", s.Name())
		}()
	}

	wg.Wait()
}

// StopAll stop all service
func (r *ServiceRegistry) StopAll() {
	var wg sync.WaitGroup
	wg.Add(len(r.services))

	for _, s := range r.services {
		go func() {
			defer wg.Done()

			log.Printf("Stopping service: '%v'...", s.Name())

			s.Stop()

			log.Printf("Stopped service: '%v'.", s.Name())
		}()
	}

	wg.Wait()
}
