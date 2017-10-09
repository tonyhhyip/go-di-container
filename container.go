package container

type container struct {
	instances map[string]interface{}
	bindings  map[string]bindBond
	alias     map[string]string
}

type bindBond struct {
	Builder BuilderFunc
	Shared  bool
}

func NewContainer() *container {
	container := container{}
	container.Flush()
	return &container
}

func (c *container) Singleton(abstract string, builder BuilderFunc) {
	c.registerBinding(abstract, builder, true)
}

func (c *container) Bind(abstract string, builder BuilderFunc) {
	c.registerBinding(abstract, builder, false)
}

func (c *container) registerBinding(abstract string, builder BuilderFunc, shared bool) {
	c.bindings[abstract] = bindBond{
		Builder: builder,
		Shared:  shared,
	}
}

func (c *container) Instance(abstract string, instance interface{}) {
	c.instances[abstract] = instance
}
func (c *container) makeWithContainer(container Container, abstract string) (instance interface{}) {
	name := c.getAlias(abstract)

	if instance, exists := c.instances[name]; exists {
		return instance
	}

	if !c.hasRegister(name) {
		return nil
	}

	builder := c.getConstructor(name)
	instance = builder(container)

	if c.isShared(name) {
		c.instances[name] = instance
	}

	return instance
}

func (c *container) Make(abstract string) (instance interface{}) {
	return c.makeWithContainer(c, abstract)
}

func (c *container) Flush() {
	c.ForgetInstances()
	c.alias = make(map[string]string)
	c.bindings = make(map[string]bindBond)
}

func (c *container) ForgetInstances() {
	c.instances = make(map[string]interface{})
}

func (c *container) ForgetInstance(abstract string) {
	delete(c.instances, abstract)
}

func (c *container) Alias(name, abstract string) {
	c.alias[name] = abstract
}

func (c *container) getAlias(name string) string {
	if val, ok := c.alias[name]; ok {
		return val
	}

	return name
}

func (c *container) getConstructor(name string) BuilderFunc {
	return c.bindings[name].Builder
}

func (c *container) isShared(name string) bool {
	return c.bindings[name].Shared
}

func (c *container) hasRegister(name string) (exists bool) {
	_, exists = c.bindings[name]
	return
}
