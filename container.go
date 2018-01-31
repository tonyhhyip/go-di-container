package container

import (
	"sync"

	"github.com/sirupsen/logrus"
)

type container struct {
	instances  *sync.Map
	bindings   *sync.Map
	alias      *sync.Map
	createLock map[string]sync.Locker
	logger     *logrus.Entry
}

type bindBond struct {
	Builder BuilderFunc
	Shared  bool
}

func NewContainer() *container {
	container := container{
		createLock: make(map[string]sync.Locker),
		logger:     logrus.New().WithField("lib", "go-di-container").WithField("component", "container"),
	}
	container.Flush()
	return &container
}

func (c *container) GetLogger() *logrus.Entry {
	return c.logger
}

func (c *container) SetDebug(debug bool) {
	if debug {
		c.logger.Logger.SetLevel(logrus.DebugLevel)
	} else {
		c.logger.Logger.SetLevel(logrus.InfoLevel)
	}
}

func (c *container) Singleton(abstract string, builder BuilderFunc) {
	c.registerBinding(abstract, builder, true)
	c.createLock[abstract] = new(sync.Mutex)
}

func (c *container) Bind(abstract string, builder BuilderFunc) {
	c.registerBinding(abstract, builder, false)
}

func (c *container) registerBinding(abstract string, builder BuilderFunc, shared bool) {
	c.logger.Debugf("Add %s (share: %t)", abstract, shared)
	c.bindings.Store(abstract, bindBond{
		Builder: builder,
		Shared:  shared,
	})
}

func (c *container) Instance(abstract string, instance interface{}) {
	c.logger.Debugf("Store instance %s", abstract)
	c.instances.Store(abstract, instance)
}

func (c *container) MakeWithContainer(container Container, abstract string) (instance interface{}) {
	name := c.getAlias(abstract)
	c.logger.Debugf("Make for %s", name)

	if instance, exists := c.instances.Load(name); exists {
		c.logger.Debugf("Load exists instance for %s", name)
		return instance
	}

	if !c.hasRegister(name) {
		c.logger.Errorf("%s is not registered", name)
		return nil
	}

	if c.isShared(name) {
		c.logger.Debugf("Lock for %s", name)
		c.createLock[name].Lock()
		defer c.createLock[name].Unlock()
	}

	builder := c.getConstructor(name)
	instance = builder(container)

	if c.isShared(name) {
		c.logger.Debugf("Store built instance for %s", name)
		c.instances.Store(name, instance)
	}

	return instance
}

func (c *container) Make(abstract string) (instance interface{}) {
	return c.MakeWithContainer(c, abstract)
}

func (c *container) Flush() {
	c.logger.Debug("Flush container")

	c.ForgetInstances()

	c.logger.Debug("Flush alias")
	c.alias = new(sync.Map)

	c.logger.Debug("Flush bindings")
	c.bindings = new(sync.Map)
}

func (c *container) ForgetInstances() {
	c.logger.Debug("Forget all instances")
	c.instances = new(sync.Map)
}

func (c *container) ForgetInstance(abstract string) {
	c.logger.Debugf("Forget instance of %s", abstract)
	c.instances.Delete(abstract)
}

func (c *container) Alias(name, abstract string) {
	c.logger.Debugf("Create alias %s => %s", name, abstract)
	c.alias.Store(name, abstract)
}

func (c *container) getAlias(name string) string {
	c.logger.Debugf("Load alias %s", name)
	if val, ok := c.alias.Load(name); ok {
		return val.(string)
	}

	return name
}

func (c *container) getConstructor(name string) BuilderFunc {
	c.logger.Debugf("Load constructor for %s", name)
	val, _ := c.bindings.Load(name)
	return val.(bindBond).Builder
}

func (c *container) isShared(name string) bool {
	c.logger.Debugf("Check %s is shared instance", name)
	val, _ := c.bindings.Load(name)
	return val.(bindBond).Shared
}

func (c *container) hasRegister(name string) (exists bool) {
	c.logger.Debugf("Check %s is registered", name)
	_, exists = c.bindings.Load(name)
	return
}
