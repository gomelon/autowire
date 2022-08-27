package example

import (
	"fmt"
)

var _ BarIface = &Bar{}

//BarIface
//grace:great
//grace:bye
type BarIface interface {
	Print()
	ToString() string
}

type Bar struct {
}

//NewDefaultBar
//+autowire.Provider
func NewDefaultBar() *Bar {
	return &Bar{}
}

func (d *Bar) Print() {
	fmt.Print("I'm BarIface. ")
}

func (d *Bar) ToString() string {
	return "I'm BarIface. "
}

type BarGreeting string

type BarAOPForGreet BarIface

var _ BarIface = &BarAOPWithGreetImpl{}

type BarAOPWithGreetImpl struct {
	BarIface
	greeting BarGreeting
}

//NewBarAOPWithGreet
//+autowire.Provider Order=5
func NewBarAOPWithGreet(Bar BarAOPForGreet, greeting BarGreeting) *BarAOPWithGreetImpl {
	return &BarAOPWithGreetImpl{greeting: greeting, BarIface: Bar}
}

func (d *BarAOPWithGreetImpl) Print() {
	fmt.Print(d.greeting)
	d.BarIface.Print()
}

type BarBye string

type BarAOPForBye BarIface

var _ BarIface = &BarAOPWithByeImpl{}

type BarAOPWithByeImpl struct {
	BarIface
	bye BarBye
}

//NewBarAOPWithBye
//+autowire.Provider Order=10
func NewBarAOPWithBye(Bar BarAOPForBye, bye BarBye) *BarAOPWithByeImpl {
	return &BarAOPWithByeImpl{bye: bye, BarIface: Bar}
}

func (d *BarAOPWithByeImpl) Print() {
	d.BarIface.Print()
	fmt.Print(d.bye)
}
