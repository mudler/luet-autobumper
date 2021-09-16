package plugins_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Luet-lab/luet-autobumper/pkg/autobumper"
	. "github.com/Luet-lab/luet-autobumper/pkg/plugins"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	copy "github.com/otiai10/copy"
	"gopkg.in/yaml.v3"
)

var _ = Describe("Add package plugin", func() {
	Context("applies", func() {
		It("detects when applying a new copy instead of inplace ", func() {
			in := AddPackage{}
			Expect(in.Apply(autobumper.LuetPackageWithLabels{
				Labels: map[string]string{"autobump.newcopy": "true"},
			})).To(BeTrue())
			Expect(in.Apply(autobumper.LuetPackageWithLabels{
				Labels: map[string]string{"autobump.newcopy": "false"},
			})).To(BeFalse())
			Expect(in.Apply(autobumper.LuetPackageWithLabels{
				Labels: map[string]string{"autobump.newcopy": ""},
			})).To(BeFalse())
			Expect(in.Apply(autobumper.LuetPackageWithLabels{
				Labels: map[string]string{},
			})).To(BeFalse())
		})

		It("bumps a package by adding a new entry instead of replacing", func() {
			dir, err := ioutil.TempDir(os.TempDir(), "")
			Expect(err).ToNot(HaveOccurred())

			err = copy.Copy("../../tests/fixtures/revdeps", dir)
			Expect(err).ToNot(HaveOccurred())

			defer os.RemoveAll(dir)

			in := AddPackage{}

			err = in.Bump(
				autobumper.LuetPackageWithLabels{
					Labels: map[string]string{
						"autobump.newcopy": "true",
					},
					LuetPackage: autobumper.LuetPackage{
						Name:     "foo",
						Path:     dir,
						Version:  "1.0",
						Category: "test",
					},
				},
				autobumper.LuetPackageWithLabels{
					LuetPackage: autobumper.LuetPackage{
						Name:     "foo",
						Path:     dir,
						Category: "test",
						Version:  "1.3",
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())

			p := &autobumper.TreeResult{}
			data, err := ioutil.ReadFile(filepath.Join(dir, "collection.yaml"))
			Expect(err).ToNot(HaveOccurred())

			err = yaml.Unmarshal(data, p)
			Expect(err).ToNot(HaveOccurred())
			Expect(p.Packages[0].Version).To(Equal("1.0"))
			Expect(p.Packages[1].Version).To(Equal("1.0"))
			Expect(p.Packages[2].Version).To(Equal("1.0+1"))
			Expect(p.Packages[3].Version).To(Equal("1.3"))
		})
	})
})
