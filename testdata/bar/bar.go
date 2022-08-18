package bar

import (
	"fmt"
	"github.com/gomelon/metas/mwire/testdata/foo"
)

var _ foo.Foo = &DefaultFoo{}

type DefaultFoo struct {
}

//NewDefaultFoo
//wire:provider
func NewDefaultFoo() *DefaultFoo {
	return &DefaultFoo{}
}

func (d *DefaultFoo) Print() {
	fmt.Print("I'm Bar. ")
}

func (d *DefaultFoo) ToString() string {
	return "I'm Bar. "
}

type Greeting string

type FooAOPForGreet foo.Foo

var _ foo.Foo = &FooAOPWithGreetImpl{}

type FooAOPWithGreetImpl struct {
	foo.Foo
	greeting Greeting
}

//NewFooAOPWithGreet
//wire:provider order=5
func NewFooAOPWithGreet(foo FooAOPForGreet, greeting Greeting) *FooAOPWithGreetImpl {
	return &FooAOPWithGreetImpl{greeting: greeting, Foo: foo}
}

func (d *FooAOPWithGreetImpl) Print() {
	fmt.Print(d.greeting)
	d.Foo.Print()
}

type Bye string

type FooAOPForBye foo.Foo

var _ foo.Foo = &FooAOPWithByeImpl{}

type FooAOPWithByeImpl struct {
	foo.Foo
	bye Bye
}

//NewFooAOPWithBye
//wire:provider order=10
func NewFooAOPWithBye(foo FooAOPForBye, bye Bye) *FooAOPWithByeImpl {
	return &FooAOPWithByeImpl{bye: bye, Foo: foo}
}

func (d *FooAOPWithByeImpl) Print() {
	d.Foo.Print()
	fmt.Print(d.bye)
}
