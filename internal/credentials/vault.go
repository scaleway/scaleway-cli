package credentials

type vault struct {
}

func (v *vault) Store(path string, key string, value string) error {
	//TODO implement me
	panic("implement me")
}

func (v *vault) RequiredEnv() []string {
	//TODO implement me
	panic("implement me")
}

func (v *vault) Resolve(path string, key string) (string, error) {
	//TODO implement me
	panic("implement me")
}
