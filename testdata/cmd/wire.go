//go:build wireinject
// +build wireinject

//The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/gomelon/metas/mwire/testdata/bar"
	"github.com/google/wire"
)

func initApp(greeting bar.Greeting, bye bar.Bye) (*App, error) {
	wire.Build(bar.ProviderSet, NewApp)
	return nil, nil
}
