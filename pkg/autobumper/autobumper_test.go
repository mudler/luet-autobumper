package autobumper_test

import (
	"github.com/Luet-lab/luet-autobumper/pkg/autobumper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeCrawler struct{ FixedVersion string }

func (f *FakeCrawler) Crawl(autobumper.LuetPackage) (bool, string) {
	return true, f.FixedVersion
}

var _ = Describe("Autobumper", func() {
	Context("Package scanning", func() {
		It("Detects packages in a tree", func() {
			ab := autobumper.New(
				autobumper.WithTreePath("../../tests/fixtures"),
				autobumper.WithCrawler(&FakeCrawler{FixedVersion: "1.99999"}),
			)
			bumps, err := ab.Run()
			Expect(err).ToNot(HaveOccurred())
			diffs := []autobumper.LuetPackage{}
			for _, d := range bumps.Diffs {
				diffs = append(diffs, d)
			}
			Expect(len(diffs)).To(Equal(1))
			Expect(diffs[0].Version).To(Equal("1.99999"))
		})

		It("Doesn't bump already existing packages in a tree", func() {
			ab := autobumper.New(
				autobumper.WithTreePath("../../tests/fixtures"),
				autobumper.WithCrawler(&FakeCrawler{FixedVersion: "1.0"}),
			)
			bumps, err := ab.Run()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(bumps.Diffs)).To(Equal(0))
		})
	})
})
