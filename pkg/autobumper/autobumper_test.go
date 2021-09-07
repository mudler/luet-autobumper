package autobumper_test

import (
	"github.com/Luet-lab/luet-autobumper/pkg/autobumper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Autobumper", func() {
	Context("Package scanning", func() {

		It("Detects files in a tree", func() {
			ab := autobumper.New()
			err := ab.Run()
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
