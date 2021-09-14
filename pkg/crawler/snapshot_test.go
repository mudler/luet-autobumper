package crawler_test

import (
	"time"

	"github.com/Luet-lab/luet-autobumper/pkg/autobumper"
	. "github.com/Luet-lab/luet-autobumper/pkg/crawler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Snapshot Crawler", func() {
	Context("current day", func() {
		It("returns default format", func() {
			s := Snapshot{}
			found, v := s.Crawl(autobumper.LuetPackageWithLabels{
				Labels: map[string]string{"autobump.strategy": "snapshot"},
			})
			Expect(found).To(BeTrue())
			Expect(v).To(Equal(time.Now().Format("20060102")))
		})
	})
})
