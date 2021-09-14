package autobumper_test

import (
	"fmt"

	catalog "github.com/Luet-lab/luet-autobumper/tests/libs"

	"github.com/Luet-lab/luet-autobumper/pkg/autobumper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
				autobumper.WithCrawler(&catalog.FakeCrawler{FixedVersion: "1.99999"}),
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
				autobumper.WithCrawler(&catalog.FakeCrawler{FixedVersion: "1.0"}),
			)
			bumps, err := ab.Run()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(bumps.Diffs)).To(Equal(0))
		})

		It("does read labels", func() {
			ab := autobumper.New(
				autobumper.WithTreePath(d),
				autobumper.WithCrawler(&catalog.FakeCrawler{VersionFromLabel: true}),
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

	Context("Plugins", func() {
		It("executes", func() {
			p := &catalog.FakePlugin{}
			ab := autobumper.New(
				autobumper.WithTreePath("../../tests/fixtures/test_tree"),
				autobumper.WithPlugin(p),
				autobumper.WithCrawler(&catalog.FakeCrawler{FixedVersion: "1.99999"}),
			)
			bumps, err := ab.Run()
			Expect(err).ToNot(HaveOccurred())
			diffs := diffsFromBumps(bumps)
			Expect(len(diffs)).To(Equal(1))
			Expect(diffs[0].Version).To(Equal("1.99999"))
			Expect(p.Bumped).To(Equal(true))
		})
	})
})
