package plugins

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Luet-lab/luet-autobumper/pkg/autobumper"
)

type Revdeps struct {
	Tree string
}

func (rdeps *Revdeps) Bump(src autobumper.LuetPackageWithLabels, dst autobumper.LuetPackageWithLabels) error {
	// Get all revdeps of the package
	// Increment the +.. field if present or add a +0
	r, err := autobumper.GetPackagesRevdeps(rdeps.Tree, fmt.Sprintf("%s/%s", src.Category, src.Name))
	if err != nil {
		return err
	}

	for _, d := range r {
		v := d.Version
		if strings.Contains(v, "+") {
			dat := strings.Split(v, "+")
			if len(dat) != 2 {
				// Skip invalid versions that we cannot bump
				continue
			}
			revVersion := dat[1]
			revVersionI, err := strconv.Atoi(revVersion)
			if err != nil {
				// we could not convert the value after '+' to int, so we skip
				continue
			}
			d.SetField("version", dat[0]+"+"+fmt.Sprint(revVersionI+1))
		} else {
			d.SetField("version", v+"+0")
		}
	}

	return nil
}

func (rdeps *Revdeps) Apply(p autobumper.LuetPackageWithLabels) bool {
	return p.Labels["autobump.revdeps"] == "true"
}
