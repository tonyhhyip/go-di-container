package container

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type kernel struct {
	Container
	providers     []ServiceProvider
	providersLock sync.Locker
	defered       *sync.Map
	bootstrap     bool
	logger        *logrus.Entry
}

func NewKernel() *kernel {
	c := NewContainer()
	kernel := kernel{
		Container:     c,
		bootstrap:     false,
		providers:     make([]ServiceProvider, 0),
		providersLock: new(sync.Mutex),
		defered:       new(sync.Map),
		logger:        c.GetLogger().WithField("component", "kernel"),
	}

	kernel.Flush()

	return &kernel
}

func (kernel *kernel) GetLogger() *logrus.Entry {
	return kernel.logger
}

func (kernel *kernel) Make(abstract string) interface{} {
	if kernel.bootstrap {
		for _, provider := range kernel.providers {
			if !provider.IsBooted() {
				provider.Boot()
			}
		}
		kernel.bootstrap = false
	}

	if _, exists := kernel.defered.Load(abstract); exists {
		kernel.loadDeferServiceProvider(abstract)
	}

	return kernel.Container.MakeWithContainer(kernel, abstract)
}

func (kernel *kernel) Flush() {
	kernel.Container.Flush()
	kernel.providersLock.Lock()
	defer kernel.providersLock.Unlock()
	kernel.GetLogger().Debugf("Clean up lock")
	kernel.providers = make([]ServiceProvider, 0)
}

func (kernel *kernel) loadDeferServiceProvider(abstract string) {
	kernel.logger.Debugf("Load Defer Service Provider %s", abstract)
	val, _ := kernel.defered.Load(abstract)
	provider := val.(ServiceProvider)
	if !provider.IsBooted() {
		provider.Register(kernel)
		boot(provider)
		kernel.appendProvider(provider)
	}

	kernel.defered.Delete(abstract)
}

func (kernel *kernel) Register(builder ServiceProviderBuilder) {
	instance := builder(kernel)

	if instance.IsDefer() {
		for _, provide := range instance.Provides() {
			kernel.defered.Store(provide, instance)
		}

		return
	}

	if !instance.IsBooted() {
		instance.Register(kernel)
		kernel.bootstrap = true

		kernel.appendProvider(instance)
	}
}

func (kernel *kernel) appendProvider(provider ServiceProvider) {
	kernel.providersLock.Lock()
	defer kernel.providersLock.Unlock()
	kernel.providers = append(kernel.providers, provider)
}
