package autobumper

type GitOptions struct {
	Auto        bool
	Signoff     bool
	CommitArgs  string
	GithubToken string
	StartBranch string
	HubArgs     string
}

type LuetOptions struct {
	PackageTreePath string
}

type Option func(cfg *Config) error

type Config struct {
	Git       *GitOptions
	Luet      *LuetOptions
	PkgApi    string
	KeepGoing bool

	// crawlers simply detect new versions from a defined source
	crawlers []crawler
	plugins  []plugin
}

// Apply applies the given options to the config, returning the first error
// encountered (if any).
func (cfg *Config) Apply(opts ...Option) error {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(cfg); err != nil {
			return err
		}
	}
	return nil
}
