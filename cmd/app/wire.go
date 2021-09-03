// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/towelong/lin-cms-go/internal"
)

func NewInjector() (*internal.Injector, error) {
	panic(wire.Build(set))
}
