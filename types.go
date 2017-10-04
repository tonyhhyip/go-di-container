package container

type BuilderFunc func(app *container) interface{}

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
