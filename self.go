package autowire

type SelfWirer interface {
	SetSelf(self any)
}
