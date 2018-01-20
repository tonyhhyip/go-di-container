package container

type BuilderFunc func(app Container) interface{}

type ServiceProviderBuilder func(app Container) ServiceProvider

type Container interface {
	makeWithContainer(container Container, abstract string) (instance interface{})

	Singleton(abstract string, builder BuilderFunc)
	Bind(abstract string, builder BuilderFunc)
	Instance(abstract string, instance interface{})
	Make(abstract string) (instance interface{})
	Flush()
	ForgetInstances()
	ForgetInstance(abstract string)
	Alias(name, abstract string)
}

type Kernel interface {
	Container

	loadDeferServiceProvider(abstract string)

	Register(builder ServiceProviderBuilder)
}

type ServiceProvider interface {
	SetContainer(container Container)
	IsDefer() bool
	IsBooted() bool
	Boot()
	Register(container Container)
	Provides() []string
}
