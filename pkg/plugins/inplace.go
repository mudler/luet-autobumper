package plugins

import (
	"github.com/Luet-lab/luet-autobumper/pkg/autobumper"
)

type Inplace struct {
}

func (inplace *Inplace) Bump(src autobumper.LuetPackageWithLabels, dst autobumper.LuetPackageWithLabels) error {
	return src.SetField("version", dst.Version)
}

func (inplace *Inplace) Apply(p autobumper.LuetPackageWithLabels) bool {
	return p.Labels["autobump.inplace"] != "false"
}
