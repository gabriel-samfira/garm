package common

type Worker interface {
	Start() error
	Stop() error
}
