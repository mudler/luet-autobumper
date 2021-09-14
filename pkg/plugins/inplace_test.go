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

var _ = Describe("Inplace plugin", func() {
	Context("applies", func() {
		It("detects when", func() {
			in := Inplace{}
			Expect(in.Apply(autobumper.LuetPackageWithLabels{
				Labels: map[string]string{"autobump.inplace": "true"},
			})).To(BeTrue())
			Expect(in.Apply(autobumper.LuetPackageWithLabels{
				Labels: map[string]string{"autobump.inplace": "false"},
			})).To(BeFalse())
			Expect(in.Apply(autobumper.LuetPackageWithLabels{
				Labels: map[string]string{"autobump.inplace": ""},
			})).To(BeTrue())
			Expect(in.Apply(autobumper.LuetPackageWithLabels{
				Labels: map[string]string{},
			})).To(BeTrue())
		})

		It("bumps a package (not a collection)", func() {
			dir, err := ioutil.TempDir(os.TempDir(), "")
			Expect(err).ToNot(HaveOccurred())

			err = copy.Copy("../../tests/fixtures/test_tree", dir)
			Expect(err).ToNot(HaveOccurred())

			defer os.RemoveAll(dir)

			in := Inplace{}

			err = in.Bump(
				autobumper.LuetPackageWithLabels{
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
						Version:  "1.3+1",
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())

			p := &autobumper.LuetPackage{}
			data, err := ioutil.ReadFile(filepath.Join(dir, "definition.yaml"))
			Expect(err).ToNot(HaveOccurred())

			err = yaml.Unmarshal(data, p)
			Expect(err).ToNot(HaveOccurred())
			Expect(p.Version).To(Equal("1.3+1"))

		})

	})
})
