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

func WithCrawler(c ...crawler) func(cfg *Config) error {
	return func(cfg *Config) error {
		cfg.crawlers = append(cfg.crawlers, c...)
		return nil
	}
}

func WithPlugin(p ...plugin) func(cfg *Config) error {
	return func(cfg *Config) error {
		cfg.plugins = append(cfg.plugins, p...)
		return nil
	}
}
