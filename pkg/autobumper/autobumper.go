package autobumper

import "github.com/hashicorp/go-multierror"

type crawler interface {
	// Crawl return a boolean and a string
	// The boolean is wether there is a new version available, and the string is the new version found
	Crawl(LuetPackageWithLabels) (bool, string)

	Apply(LuetPackageWithLabels) bool
}

type plugin interface {
	Apply(LuetPackageWithLabels) bool
	Bump(LuetPackageWithLabels, LuetPackageWithLabels) error
}

type AutoBumper struct {
	config Config
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

// TODO: Retrieve labels for behavior. Decide how to bump (inplace, add a new package)

func (ab *AutoBumper) Bump(src LuetPackage, bumps Bumps) error {
	for _, p := range ab.config.plugins {
		if p.Apply(src.WithLabels()) {
			if err := p.Bump(src.WithLabels(), bumps.Diffs[src].WithLabels()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (ab *AutoBumper) Run() (Bumps, error) {

	b := Bumps{Diffs: make(map[LuetPackage]LuetPackage)}

	packs, err := GetPackages(ab.config.Luet.PackageTreePath)
	if err != nil {
		return b, err
	}

	// TODO: Collect error instead of returning immediately
	// TOO: crowlers retrieve labels for behavior
	for _, p := range packs {
		for _, c := range ab.config.crawlers {
			pLabels := p.WithLabels()
			if c.Apply(pLabels) {
				if found, version := c.Crawl(pLabels); found && !Packages(packs).In(p.WithVersion(version)) {
					b.Diffs[p] = p.WithVersion(version)
					if berr := ab.Bump(p, b); berr != nil {
						err = multierror.Append(err, berr)
					}
				}
			}
		}
	}

	return b, err
}
