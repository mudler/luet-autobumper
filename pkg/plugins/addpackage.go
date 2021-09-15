package plugins

import (
	"fmt"
	"path/filepath"

	"github.com/Luet-lab/luet-autobumper/pkg/autobumper"
	"github.com/Luet-lab/luet-autobumper/pkg/utils"
	"github.com/otiai10/copy"
)

type AddPackage struct {}

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
		return f.SetField("version", dst.Version)
	}

	// Read the collection, and find the src
	coll, err := autobumper.ReadCollection(src.Path)
	if err != nil {
		return err
	}

	index, err := autobumper.Collection(coll).Find(src.WithLabels())
	if err != nil {
		return err
	}

	// copy from src to a new list item with yq:
	// We use yq to preserve fields outside of structs,
	// comments, etc.
	if err := utils.YQ(
		fmt.Sprintf(".packages[%d] = .packages[%d]",
			len(coll),
			index,
		),
		filepath.Join(src.Path, "collection.yaml"),
	); err != nil {
		return err
	}

	// set new version
	if err := utils.YQ(
		fmt.Sprintf(".packages[%d].version = %s",
			len(coll),
			dst.Version,
		),
		filepath.Join(src.Path, "collection.yaml"),
	); err != nil {
		return err
	}

	return nil
}

func (add *AddPackage) Apply(p autobumper.LuetPackageWithLabels) bool {
	return p.Labels["autobump.newcopy"] == "true" || p.Labels["autobump.inplace"] == "false"
}
