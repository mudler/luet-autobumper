package crawler

import (
	"time"

	"github.com/Luet-lab/luet-autobumper/pkg/autobumper"
)

type Snapshot struct {
	Format string
}

func (s *Snapshot) Apply(p autobumper.LuetPackageWithLabels) bool {
	return p.Strategy() == "snapshot"
}

func (s *Snapshot) Crawl(p autobumper.LuetPackageWithLabels) (bool, string) {
	format := "20060102"
	if s.Format != "" {
		format = s.Format
	}
	currentTime := time.Now()
	return true, currentTime.Format(format)
}
