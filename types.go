package container

type BuilderFunc func(app *container) interface{}

type ServiceProviderBuilder func(container Container) ServiceProvider

type Container interface {
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
