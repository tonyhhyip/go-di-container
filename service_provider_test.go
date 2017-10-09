package container

import (
	"fmt"
	"testing"
)

type aServiceProvider struct {
	*AbstractServiceProvider
}

func (*aServiceProvider) Provides() []string {
	return []string{
		"foo",
	}
}

func (*aServiceProvider) Register(app Container) {
	app.Bind("foo", func(app Container) interface{} {
		return "a"
	})
}

type bServiceProvider struct {
	*AbstractServiceProvider
}

func (*bServiceProvider) Provides() []string {
	return []string{
		"bar",
	}
}

func (*bServiceProvider) Register(app Container) {
	app.Bind("bar", func(app Container) interface{} {
		a := app.Make("foo").(string)
		return fmt.Sprintf("%sb", a)
	})
}

func TestMultipleDeferRegister(t *testing.T) {
	kernel := NewKernel()
	kernel.Register(func(app Container) ServiceProvider {
		sp := aServiceProvider{
			NewAbstractServiceProvider(true),
		}
		sp.SetContainer(app)
		return &sp
	})
	kernel.Register(func(app Container) ServiceProvider {
		sp := bServiceProvider{
			NewAbstractServiceProvider(true),
		}
		sp.SetContainer(app)
		return &sp
	})

	b := kernel.Make("bar").(string)
	if b != "ab" {
		t.Error("Wrong output")
	}
}
