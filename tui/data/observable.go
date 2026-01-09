package data

type Observable interface {
	Observe(func())
}
