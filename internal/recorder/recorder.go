package recorder

type Recorder interface {
	StartStream() error
	StopStream()
}
