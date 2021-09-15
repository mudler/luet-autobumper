package plugins

import (
	"github.com/Luet-lab/luet-autobumper/pkg/autobumper"
	"gopkg.in/op/go-logging.v1"
)

type Inplace struct {
}

var log = logging.MustGetLogger("plugins")

func (inplace *Inplace) Bump(src autobumper.LuetPackageWithLabels, dst autobumper.LuetPackageWithLabels) error {
	log.Info("Bumping %s to %s", src, dst)
	return src.SetField("version", dst.Version)
}

func (inplace *Inplace) Apply(p autobumper.LuetPackageWithLabels) bool {
	return p.Labels["autobump.inplace"] != "false"
}
