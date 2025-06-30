package authstrategy

type Strategy interface {
	Execute() (string, error)
}

func StrategyFactory(conf map[string]string) Strategy {
	switch conf["name"] {
	case "gigachat":
		return GigachatStrategy{Conf: conf}
	default:
		return DefaultStrategy{Conf: conf}
	}
}
