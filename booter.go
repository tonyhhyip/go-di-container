package container

type canBoot interface {
	Boot()
}

func boot(obj interface{}) {
	if booter, ok := obj.(canBoot); ok {
		booter.Boot()
	}
}
