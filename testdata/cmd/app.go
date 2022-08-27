package main

import (
	"github.com/gomelon/autowire/testdata/example"
)

type App struct {
	foo  example.FooIface
	bar  example.BarIface
	aMao example.AMaoIface
}

func NewApp(foo example.FooIface, bar example.BarIface, aMao example.AMaoIface) *App {
	return &App{foo: foo, bar: bar, aMao: aMao}
}
