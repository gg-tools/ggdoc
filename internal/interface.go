package internal

type Titled interface {
	GetTitle() (string, error)
}
