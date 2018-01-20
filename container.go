package container

import "sync"

type container struct {
	instances  *sync.Map
	bindings   *sync.Map
	alias      *sync.Map
	createLock map[string]sync.Locker
}

type bindBond struct {
	Builder BuilderFunc
	Shared  bool
}

func NewContainer() *container {
	container := container{
		createLock: make(map[string]sync.Locker),
	}
	container.Flush()
	return &container
}

func (c *container) Singleton(abstract string, builder BuilderFunc) {
	c.registerBinding(abstract, builder, true)
	c.createLock[abstract] = new(sync.Mutex)
}

func (c *container) Bind(abstract string, builder BuilderFunc) {
	c.registerBinding(abstract, builder, false)
}

func (c *container) registerBinding(abstract string, builder BuilderFunc, shared bool) {
	c.bindings.Store(abstract, bindBond{
		Builder: builder,
		Shared:  shared,
	})
}

func (c *container) Instance(abstract string, instance interface{}) {
	c.instances.Store(abstract, instance)
}

func (c *container) MakeWithContainer(container Container, abstract string) (instance interface{}) {
	name := c.getAlias(abstract)

	if instance, exists := c.instances.Load(name); exists {
		return instance
	}

	if !c.hasRegister(name) {
		return nil
	}

	if c.isShared(name) {
		c.createLock[name].Lock()
		defer c.createLock[name].Unlock()
	}

	builder := c.getConstructor(name)
	instance = builder(container)

	if c.isShared(name) {
		c.instances.Store(name, instance)
	}

	return instance
}

func (c *container) Make(abstract string) (instance interface{}) {
	return c.MakeWithContainer(c, abstract)
}

func (c *container) Flush() {
	c.ForgetInstances()
	c.alias = new(sync.Map)
	c.bindings = new(sync.Map)
}

func (c *container) ForgetInstances() {
	c.instances = new(sync.Map)
}

func (c *container) ForgetInstance(abstract string) {
	c.instances.Delete(abstract)
}

func (c *container) Alias(name, abstract string) {
	c.alias.Store(name, abstract)
}

func (c *container) getAlias(name string) string {
	if val, ok := c.alias.Load(name); ok {
		return val.(string)
	}

	return name
}

func (c *container) getConstructor(name string) BuilderFunc {
	val, _ := c.bindings.Load(name)
	return val.(bindBond).Builder
}

func (c *container) isShared(name string) bool {
	val, _ := c.bindings.Load(name)
	return val.(bindBond).Shared
}

func (c *container) hasRegister(name string) (exists bool) {
	_, exists = c.bindings.Load(name)
	return
}
