package autobumper_test

import (
	"fmt"

	"github.com/Luet-lab/luet-autobumper/pkg/autobumper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeCrawler struct {
	FixedVersion     string
	VersionFromLabel bool
}

func (f *FakeCrawler) Crawl(p autobumper.LuetPackage) (bool, string) {
	if f.VersionFromLabel {
		labels, err := p.ReadLabels()
		Expect(err).ToNot(HaveOccurred())
		if v, ok := labels["version"]; ok {
			return true, v
		} else {
			return false, ""
		}
	}

	return true, f.FixedVersion
}

func diffsFromBumps(bumps autobumper.Bumps) []autobumper.LuetPackage {
	diffs := []autobumper.LuetPackage{}
	for _, d := range bumps.Diffs {
		diffs = append(diffs, d)
	}
	return diffs
}

func smokeTest(d string) {
	Context(fmt.Sprintf("package scanning in '%s'", d), func() {

		It("detects packages in a tree", func() {
			ab := autobumper.New(
				autobumper.WithTreePath(d),
				autobumper.WithCrawler(&FakeCrawler{FixedVersion: "1.99999"}),
			)
			bumps, err := ab.Run()
			Expect(err).ToNot(HaveOccurred())
			diffs := diffsFromBumps(bumps)
			Expect(len(diffs)).To(Equal(1))
			Expect(diffs[0].Version).To(Equal("1.99999"))
		})

		It("doesn't bump already existing packages in a tree", func() {
			ab := autobumper.New(
				autobumper.WithTreePath(d),
				autobumper.WithCrawler(&FakeCrawler{FixedVersion: "1.0"}),
			)
			bumps, err := ab.Run()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(bumps.Diffs)).To(Equal(0))
		})

		It("does read labels", func() {
			ab := autobumper.New(
				autobumper.WithTreePath(d),
				autobumper.WithCrawler(&FakeCrawler{VersionFromLabel: true}),
			)
			bumps, err := ab.Run()
			Expect(err).ToNot(HaveOccurred())
			diffs := diffsFromBumps(bumps)
			Expect(len(diffs)).To(Equal(1))
			Expect(diffs[0].Version).To(Equal("baz"))
		})
	})
}

var _ = Describe("Autobumper", func() {
	smokeTest("../../tests/fixtures/test_tree")
	smokeTest("../../tests/fixtures/collection")
})
