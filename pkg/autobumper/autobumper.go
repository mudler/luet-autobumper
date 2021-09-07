package autobumper

type crawler interface {
	Crawl(LuetPackage) (bool, string)
}

type AutoBumper struct {
	config   Config
	crawlers []crawler
}

func New(p ...Option) *AutoBumper {
	c := Config{}
	c.Apply(p...)

	return &AutoBumper{
		config: c,
	}
}

// TODO: maybe use interfaces here too?
func (ab *AutoBumper) Bump(p LuetPackage, v string) error {
	return nil
}

func (ab *AutoBumper) Run() error {
	packs, err := ab.getPackages(ab.config.Luet.PackageTreePath)
	if err != nil {
		return err
	}

	// TODO: Collect error instead of returning immediately
	for _, p := range packs {
		for _, c := range ab.crawlers {
			if found, version := c.Crawl(p); found {
				if err := ab.Bump(p, version); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
