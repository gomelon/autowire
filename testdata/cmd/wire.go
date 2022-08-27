//go:build wireinject
// +build wireinject

//The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/gomelon/autowire/testdata/example"
	"github.com/google/wire"
)

func initApp(barGreeting example.BarGreeting, barBye example.BarBye,
	fooGreeting example.FooGreeting, fooBye example.FooBye) (*App, error) {
	wire.Build(example.ProviderSet, NewApp)
	return nil, nil
}
