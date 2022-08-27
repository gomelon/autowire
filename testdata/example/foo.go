package example

import (
	"fmt"
)

var _ FooIface = &Foo{}

//FooIface
//grace:great
//grace:bye
type FooIface interface {
	Print()
	ToString() string
}

type Foo struct {
}

//NewDefaultFoo
//+autowire.Provider
func NewDefaultFoo() *Foo {
	return &Foo{}
}

func (d *Foo) Print() {
	fmt.Print("I'm FooIface. ")
}

func (d *Foo) ToString() string {
	return "I'm FooIface. "
}

type FooGreeting string

type FooAOPForGreet FooIface

var _ FooIface = &FooAOPWithGreetImpl{}

type FooAOPWithGreetImpl struct {
	FooIface
	greeting FooGreeting
}

//NewFooAOPWithGreet
//+autowire.Provider Order=5
func NewFooAOPWithGreet(foo FooAOPForGreet, greeting FooGreeting) *FooAOPWithGreetImpl {
	return &FooAOPWithGreetImpl{greeting: greeting, FooIface: foo}
}

func (d *FooAOPWithGreetImpl) Print() {
	fmt.Print(d.greeting)
	d.FooIface.Print()
}

type FooBye string

type FooAOPForBye FooIface

var _ FooIface = &FooAOPWithByeImpl{}

type FooAOPWithByeImpl struct {
	FooIface
	bye FooBye
}

//NewFooAOPWithBye
//+autowire.Provider Order=10
func NewFooAOPWithBye(foo FooAOPForBye, bye FooBye) *FooAOPWithByeImpl {
	return &FooAOPWithByeImpl{bye: bye, FooIface: foo}
}

func (d *FooAOPWithByeImpl) Print() {
	d.FooIface.Print()
	fmt.Print(d.bye)
}
