package container

import "sync"

type syncKernel struct {
	*kernel
	deferedLock   sync.Locker
	providersLock sync.Locker
}

func NewSyncKernel() *syncKernel {
	base := NewKernel()

	return &syncKernel{
		kernel:        base,
		deferedLock:   &sync.Mutex{},
		providersLock: &sync.Mutex{},
	}
}

func (kernel *syncKernel) Flush() {
	kernel.providersLock.Lock()
	defer kernel.providersLock.Unlock()
	kernel.kernel.Flush()
}

func (kernel *syncKernel) Make(abstract string) interface{} {
	kernel.providersLock.Lock()
	defer kernel.providersLock.Unlock()

	kernel.deferedLock.Lock()
	defer kernel.deferedLock.Unlock()

	return kernel.kernel.Make(abstract)
}

func (kernel *syncKernel) Register(builder ServiceProviderBuilder) {
	kernel.providersLock.Lock()
	defer kernel.providersLock.Unlock()

	kernel.deferedLock.Lock()
	defer kernel.deferedLock.Unlock()

	kernel.kernel.Register(builder)
}
