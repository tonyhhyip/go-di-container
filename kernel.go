package container

type kernel struct {
	Container
	providers []ServiceProvider
	defered   map[string]ServiceProvider
	bootstrap bool
}

func NewKernel() *kernel {
	kernel := kernel{
		Container: NewContainer(),
		bootstrap: false,
		providers: make([]ServiceProvider, 0),
		defered:   make(map[string]ServiceProvider),
	}

	kernel.Flush()

	return &kernel
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

	if _, exists := kernel.defered[abstract]; exists {
		kernel.loadDeferServiceProvider(abstract)
	}

	return kernel.Container.MakeWithContainer(kernel, abstract)
}

func (kernel *kernel) Flush() {
	kernel.Container.Flush()
	kernel.providers = make([]ServiceProvider, 0)
}

func (kernel *kernel) loadDeferServiceProvider(abstract string) {
	provider := kernel.defered[abstract]
	if !provider.IsBooted() {
		provider.Register(kernel)
		boot(provider)
		kernel.providers = append(kernel.providers, provider)
	}

	delete(kernel.defered, abstract)
}

func (kernel *kernel) Register(builder ServiceProviderBuilder) {
	instance := builder(kernel)

	if instance.IsDefer() {
		for _, provide := range instance.Provides() {
			kernel.defered[provide] = instance
		}

		return
	}

	if !instance.IsBooted() {
		instance.Register(kernel)
		kernel.bootstrap = true
		kernel.providers = append(kernel.providers, instance)
	}
}
