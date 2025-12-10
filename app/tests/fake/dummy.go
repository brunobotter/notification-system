package fake

type DummyWriteCloser struct{}

func (d *DummyWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}

func (d *DummyWriteCloser) Close() error {
	return nil
}
