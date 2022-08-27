package example

import "fmt"

var _ AMaoIface = &AMao{}

type AMaoIface interface {
	Print()
	Name() string
	SetSelf(self any)
}

type AMao struct {
	self AMaoIface
}

//NewAMao
//+autowire.Provider
func NewAMao() *AMao {
	return &AMao{}
}

func (a *AMao) SetSelf(self any) {
	a.self = self.(AMaoIface)
}

func (a *AMao) Print() {
	fmt.Println(a.self.Name())
}

func (a *AMao) Name() string {
	return "a mao"
}

type AMaoForGreet AMaoIface

var _ AMaoIface = &AMaoWithGreetImpl{}

type AMaoWithGreetImpl struct {
	AMaoIface
}

//NewAMaoWithGreetImpl
//+autowire.Provider
func NewAMaoWithGreetImpl(aMao AMaoForGreet) *AMaoWithGreetImpl {
	return &AMaoWithGreetImpl{AMaoIface: aMao}
}

func (a *AMaoWithGreetImpl) Name() string {
	return "Hi, I'm " + a.AMaoIface.Name()
}
