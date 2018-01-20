package container

import "sync"

type syncContainer struct {
	*container
	aliasLock     sync.Locker
	bindingLock   sync.Locker
	instancesLock sync.Locker
}

func NewSyncContainer() *syncContainer {
	return &syncContainer{
		container:     NewContainer(),
		aliasLock:     &sync.Mutex{},
		bindingLock:   &sync.Mutex{},
		instancesLock: &sync.Mutex{},
	}
}

func (c *syncContainer) registerBinding(abstract string, builder BuilderFunc, shared bool) {
	c.bindingLock.Lock()
	defer c.bindingLock.Unlock()
	c.container.registerBinding(abstract, builder, shared)
}

func (c *syncContainer) Instance(abstract string, instance interface{}) {
	c.instancesLock.Lock()
	defer c.instancesLock.Unlock()
	c.container.Instance(abstract, instance)
}

func (c *syncContainer) Flush() {
	c.aliasLock.Lock()
	defer c.aliasLock.Unlock()
	c.bindingLock.Lock()
	defer c.bindingLock.Unlock()
	c.container.Flush()
}

func (c *syncContainer) ForgetInstances() {
	c.instancesLock.Lock()
	defer c.instancesLock.Unlock()
	c.instances = make(map[string]interface{})
}

func (c *syncContainer) ForgetInstance(abstract string) {
	c.instancesLock.Lock()
	defer c.instancesLock.Unlock()
	delete(c.instances, abstract)
}

func (c *syncContainer) Alias(name, abstract string) {
	c.aliasLock.Lock()
	defer c.aliasLock.Unlock()
	c.alias[name] = abstract
}
