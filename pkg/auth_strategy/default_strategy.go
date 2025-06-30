package authstrategy

type DefaultStrategy struct {
	Conf map[string]string
}

func (s DefaultStrategy) Execute() (string, error) {
	return s.Conf["key"], nil
}
