package autobumper_test

import (
	"fmt"

	"github.com/Luet-lab/luet-autobumper/pkg/autobumper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeCrawler struct{}

func (f *FakeCrawler) Crawl(autobumper.LuetPackage) (bool, string) {
	return true, "1.99999"
}

var _ = Describe("Autobumper", func() {
	Context("Package scanning", func() {

		It("Detects files in a tree", func() {
			ab := autobumper.New(
				autobumper.WithTreePath("../../tests/fixtures"),
				autobumper.WithCrawler(&FakeCrawler{}),
			)
			bumps, err := ab.Run()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(bumps.Diffs)).ToNot(Equal(0))
			fmt.Println(bumps.Diffs)
		})
	})
})
