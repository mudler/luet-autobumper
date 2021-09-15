package plugins

import (
	"fmt"
	"path/filepath"

	"github.com/Luet-lab/luet-autobumper/pkg/autobumper"
	"github.com/otiai10/copy"
)

type AddPackage struct {
}

func (add *AddPackage) Bump(src autobumper.LuetPackageWithLabels, dst autobumper.LuetPackageWithLabels) error {
	log.Info("Bumping %s to %s", src, dst)

	if !src.IsCollection() {
		// If it's not a collection, copy the spec to a dir named
		// after "$PACKAGE_NAME-$PACKAGE_VERSION" on the same top level folder
		// where we have found the spec
		path := src.GetPath()
		dir := filepath.Base(path)
		err := copy.Copy(path, filepath.Join(dir, fmt.Sprintf("%s-%s", dst.Name, dst.Version)))
		if err != nil {
			return err
		}
		f := &dst
		f.Path = filepath.Join(dir, fmt.Sprintf("%s-%s", dst.Name, dst.Version))
		f.SetField("version", dst.Version)
	} else {
		// TODO:
		// Copy src portion over and append
		// modify version

	}
	return src.SetField("version", dst.Version)
}

func (add *AddPackage) Apply(p autobumper.LuetPackageWithLabels) bool {
	if p.Labels["autobump.newcopy"] != "true" {
		return false
	}

	return true
}
