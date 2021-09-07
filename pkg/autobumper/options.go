package autobumper

func AutoGit(b bool) func(cfg *Config) error {
	return func(cfg *Config) error {
		cfg.Git.Auto = b
		return nil
	}
}

func WithTreePath(t string) func(cfg *Config) error {
	return func(cfg *Config) error {
		cfg.Luet.PackageTreePath = t
		return nil
	}
}
