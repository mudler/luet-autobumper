package catalog

import "github.com/Luet-lab/luet-autobumper/pkg/autobumper"

type FakeCrawler struct {
	FixedVersion     string
	VersionFromLabel bool
}

func (f *FakeCrawler) Apply(autobumper.LuetPackageWithLabels) bool { return true }
func (f *FakeCrawler) Crawl(p autobumper.LuetPackageWithLabels) (bool, string) {
	if f.VersionFromLabel {
		labels := p.Labels
		if v, ok := labels["version"]; ok {
			return true, v
		} else {
			return false, ""
		}
	}

	return true, f.FixedVersion
}

type FakePlugin struct {
	Bumped bool
}

func (p *FakePlugin) Apply(autobumper.LuetPackageWithLabels) bool {
	return true
}

func (p *FakePlugin) Bump(autobumper.LuetPackageWithLabels, autobumper.LuetPackageWithLabels) error {
	p.Bumped = true
	return nil
}
