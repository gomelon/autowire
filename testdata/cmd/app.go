package main

import (
	"github.com/gomelon/autowire/testdata/foo"
)

type App struct {
	foo foo.Foo
}

func NewApp(foo foo.Foo) *App {
	return &App{foo: foo}
}
