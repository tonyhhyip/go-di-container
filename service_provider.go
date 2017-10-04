package container

type AbstractServiceProvider struct {
	container Container
	defered   bool
	booted    bool
}

func NewAbstractServiceProvider(deferLoading bool) *AbstractServiceProvider {
	return &AbstractServiceProvider{
		defered: deferLoading,
		booted:  false,
	}
}

func (sp *AbstractServiceProvider) IsBooted() bool {
	return sp.booted
}

func (sp *AbstractServiceProvider) Boot() {
	sp.booted = true
}

func (sp *AbstractServiceProvider) IsDefer() bool {
	return sp.defered
}

func (sp *AbstractServiceProvider) Provides() []string {
	return []string{}
}

func (sp *AbstractServiceProvider) SetContainer(container Container) {
	sp.container = container
}
