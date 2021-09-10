package autobumper

import (
	"encoding/json"
	"fmt"

	"github.com/Luet-lab/luet-autobumper/pkg/utils"
	"github.com/pkg/errors"
)

type treeResult struct {
	Packages []LuetPackage `json:"packages"`
}
type LuetPackage struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Category string `json:"category"`
	Version  string `json:"version"`
}

func (p LuetPackage) WithVersion(v string) LuetPackage {
	pp := &p
	pp.Version = v
	return *pp
}

type Packages []LuetPackage

func (p Packages) In(LuetPackage) bool {
	// TODO: Check if given package @ version is there in the set
	return false
}

func (ab *AutoBumper) getPackages(dir string) ([]LuetPackage, error) {
	jsonPacks, err := utils.RunSH(fmt.Sprintf("luet tree pkglist --tree %s -o json", ab.config.Luet.PackageTreePath))
	if err != nil {
		return []LuetPackage{}, errors.Wrap(err, "failed getting packages with luet")
	}

	packages := &treeResult{}
	if err := json.Unmarshal([]byte(jsonPacks), packages); err != nil {
		return []LuetPackage{}, errors.Wrap(err, "failed getting packages with luet")
	}

	return packages.Packages, nil
}
