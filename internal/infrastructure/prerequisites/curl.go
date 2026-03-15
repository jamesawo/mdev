package prerequisites

type Curl struct{}

func (c *Curl) Name() string {
	return "curl"
}

func (c *Curl) Check() bool {
	return true
}

func (c *Curl) Install() error {
	return nil
}

func init() {
	Register(&Curl{})
}
