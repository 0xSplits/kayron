package deploy

type Interface interface {
	Empty() bool
	Verify() error
}
