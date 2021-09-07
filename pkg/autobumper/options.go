package autobumper

func AutoGit(b bool) func(cfg *Config) error {
	return func(cfg *Config) error {
		cfg.Git.Auto = b
		return nil
	}
}
