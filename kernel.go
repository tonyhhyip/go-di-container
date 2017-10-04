package container

type Kernel struct {
	*container
	providers []ServiceProvider
	defered   map[string]ServiceProvider
	bootstrap bool
}

func NewKernel() *Kernel {
	kernel := Kernel{
		container: NewContainer(),
		bootstrap: false,
		providers: make([]ServiceProvider, 0),
		defered:   make(map[string]ServiceProvider),
	}
	return &kernel
}

func (kernel *Kernel) Boot() {

}

func (kernel *Kernel) Make(abstract string) interface{} {
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

	return kernel.container.Make(abstract)
}

func (kernel *Kernel) Flush() {
	kernel.container.Flush()
	kernel.providers = make([]ServiceProvider, 0)
}

func (kernel *Kernel) loadDeferServiceProvider(abstract string) {
	provider := kernel.defered[abstract]
	if !provider.IsBooted() {
		provider.Register(kernel)
		boot(provider)
		kernel.providers = append(kernel.providers, provider)
	}

	delete(kernel.defered, abstract)
}

func (kernel *Kernel) Register(builder ServiceProviderBuilder) {
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
