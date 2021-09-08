package autobumper

type crawler interface {
	Crawl(LuetPackage) (bool, string)
}

type AutoBumper struct {
	config   Config
	crawlers []crawler
}

type Bumps struct {
	Diffs map[LuetPackage]LuetPackage
}

func New(p ...Option) *AutoBumper {
	c := Config{
		Git:  &GitOptions{},
		Luet: &LuetOptions{},
	}
	c.Apply(p...)

	return &AutoBumper{
		config: c,
	}
}

// TODO: maybe use interfaces here too?
func (ab *AutoBumper) Bump(p LuetPackage, v string) error {
	// TODO: Retrieve labels for behavior. Decide how to bump (inplace, add a new package)
	return nil
}

func (ab *AutoBumper) Run() (Bumps, error) {

	b := Bumps{}

	packs, err := ab.getPackages(ab.config.Luet.PackageTreePath)
	if err != nil {
		return b, err
	}

	// TODO: Collect error instead of returning immediately
	// TOO: crowlers retrieve labels for behavior

	for _, p := range packs {
		for _, c := range ab.crawlers {
			if found, version := c.Crawl(p); found && !Packages(packs).In(p.Version(version)) {
				if err := ab.Bump(p, version); err != nil {
					return b, err
				}
			}
		}
	}
	return b, nil
}
