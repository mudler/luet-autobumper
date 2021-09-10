package autobumper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/Luet-lab/luet-autobumper/pkg/utils"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	collectionFile = "collection.yaml"
	definitionFile = "definition.yaml"
)

type treeResult struct {
	Packages []LuetPackage `json:"packages"`
}
type LuetPackage struct {
	Name     string `json:"name" yaml:"name"`
	Path     string `json:"path"`
	Category string `json:"category" yaml:"category"`
	Version  string `json:"version" yaml:"version"`
}

type LuetPackageWithLabels struct {
	LuetPackage `yaml:",inline"`
	Labels      map[string]string `yaml:"labels"`
}

// IsCollection returns true if the package is part of a collection
func (p LuetPackage) IsCollection() bool {
	return utils.Exists(filepath.Join(p.Path, collectionFile))
}

type packagesLabels struct {
	Packages []LuetPackageWithLabels `yaml:"packages"`
}

func (p LuetPackage) Match(pp LuetPackage) bool {
	return p.Category == pp.Category &&
		p.Name == pp.Name &&
		p.Version == pp.Version
}

func (p LuetPackage) ReadLabels() (map[string]string, error) {
	result := map[string]string{}
	// If it is a collection, we have to loop over and check which one is the corresponding package we are looking into
	if p.IsCollection() {
		res := &packagesLabels{}
		dat, err := ioutil.ReadFile(filepath.Join(p.Path, collectionFile))
		if err != nil {
			return result, err
		}
		if err := yaml.Unmarshal(dat, res); err != nil {
			return result, err
		}
		for _, ps := range res.Packages {
			if ps.Match(p) {
				result = ps.Labels
				break
			}
		}
	} else {
		// If it's not a collection, we simply parse the package and return the labels
		res := &LuetPackageWithLabels{}

		dat, err := ioutil.ReadFile(filepath.Join(p.Path, definitionFile))
		if err != nil {
			return result, err
		}
		if err := yaml.Unmarshal(dat, res); err != nil {
			return result, err
		}
		result = res.Labels
	}

	return result, nil
}

// WithVersion returns a new package copy with the version changed
func (p LuetPackage) WithVersion(v string) LuetPackage {
	pp := &p
	pp.Version = v
	return *pp
}

type Packages []LuetPackage

// In checks if the package is contained in the slice or not
func (pp Packages) In(c LuetPackage) bool {
	for _, p := range pp {
		if p.Version == c.Version {
			return true
		}
	}
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
