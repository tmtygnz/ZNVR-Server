package cameras

type Camera interface {
	GetStreamUrl() string
	GetType() string
	GetName() string
}
