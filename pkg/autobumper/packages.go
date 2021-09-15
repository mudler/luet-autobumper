package autobumper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Luet-lab/luet-autobumper/pkg/utils"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"github.com/pkg/errors"
	logging "gopkg.in/op/go-logging.v1"
	"gopkg.in/yaml.v2"
)

const (
	collectionFile = "collection.yaml"
	definitionFile = "definition.yaml"
)

type TreeResult struct {
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

type Collection []LuetPackageWithLabels

// IsCollection returns true if the package is part of a collection
func (p LuetPackage) IsCollection() bool {
	return utils.Exists(filepath.Join(p.Path, collectionFile))
}

func (p LuetPackage) Label(s string) string {
	labels, _ := p.ReadLabels()
	if v, ok := labels[s]; ok {
		return v
	}
	return ""
}

func (p LuetPackage) SetField(field, value string) error {
	// Mute logging from yqlib to just errors
	var f = logging.MustStringFormatter(
		`%{color}%{time:15:04:05} %{shortfunc} [%{level:.4s}]%{color:reset} %{message}`,
	)
	var backend = logging.AddModuleLevel(
		logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), f))

	backend.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend)

	var completedSuccessfully bool
	var expr string
	if !p.IsCollection() {

		expr = fmt.Sprintf(".%s = \"%s\"", field, value)
	} else {
		// 1: find the package inside the collection (index)
		coll, err := ReadCollection(p.Path)
		if err != nil {
			return err
		}

		index, err := Collection(coll).Find(p.WithLabels())
		if err != nil {
			return err
		}
		// 2: find the respective yaml.Node
		expr = fmt.Sprintf(".packages[%d].%s = \"%s\"", index, field, value)

	}
	format, err := yqlib.OutputFormatFromString("yaml")
	if err != nil {
		return err
	}
	writeInPlaceHandler := yqlib.NewWriteInPlaceHandler(p.GetPath())
	out, err := writeInPlaceHandler.CreateTempFile()
	if err != nil {
		return err
	}
	// need to indirectly call the function so  that completedSuccessfully is
	// passed when we finish execution as opposed to now
	defer func() { writeInPlaceHandler.FinishWriteInPlace(completedSuccessfully) }()

	printer := yqlib.NewPrinter(out, format, false, false, 0, false)

	streamEvaluator := yqlib.NewStreamEvaluator()

	err = streamEvaluator.EvaluateFiles(expr, []string{p.GetPath()}, printer, true)
	completedSuccessfully = err == nil
	return err
}

func (p LuetPackage) GetPath() string {
	if p.IsCollection() {
		return filepath.Join(p.Path, collectionFile)
	}
	return filepath.Join(p.Path, definitionFile)
}

func (p LuetPackage) WithLabels() LuetPackageWithLabels {
	labels, _ := p.ReadLabels()
	return LuetPackageWithLabels{LuetPackage: p, Labels: labels}
}

func (p LuetPackage) Strategy() string {
	return strings.ToLower(p.Label("autobump.strategy"))
}

func (p LuetPackageWithLabels) Strategy() string {
	return strings.ToLower(p.Labels["autobump.strategy"])
}

type packagesLabels struct {
	Packages []LuetPackageWithLabels `yaml:"packages"`
}

func (p LuetPackage) Match(pp LuetPackage) bool {
	return p.Category == pp.Category &&
		p.Name == pp.Name &&
		p.Version == pp.Version
}

func ReadCollection(src string) ([]LuetPackageWithLabels, error) {
	res := &packagesLabels{}
	if !strings.HasSuffix(src, collectionFile) {
		src = filepath.Join(src, collectionFile)
	}
	dat, err := ioutil.ReadFile(src)
	if err != nil {
		return []LuetPackageWithLabels{}, err
	}
	if err := yaml.Unmarshal(dat, res); err != nil {
		return []LuetPackageWithLabels{}, err
	}
	return res.Packages, nil
}

func (c Collection) Find(p LuetPackageWithLabels) (int, error) {
	for i, ps := range c {
		if ps.Match(p.LuetPackage) {
			return i, nil
		}
	}
	return 0, errors.New("not found")
}

func (p LuetPackage) ReadLabels() (map[string]string, error) {
	result := map[string]string{}
	// If it is a collection, we have to loop over and check which one is the corresponding package we are looking into
	if p.IsCollection() {
		coll, err := ReadCollection(p.Path)
		if err != nil {
			return result, err
		}
		for _, ps := range coll {
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

func GetPackages(dir string) ([]LuetPackage, error) {
	jsonPacks, err := utils.RunSH(fmt.Sprintf("luet tree pkglist --tree %s -o json", dir))
	if err != nil {
		return []LuetPackage{}, errors.Wrap(err, "failed getting packages with luet")
	}

	packages := &TreeResult{}
	if err := json.Unmarshal([]byte(jsonPacks), packages); err != nil {
		return []LuetPackage{}, errors.Wrap(err, "failed getting packages with luet")
	}

	return packages.Packages, nil
}

func GetPackagesRevdeps(dir, revdep string) ([]LuetPackage, error) {
	jsonPacks, err := utils.RunSH(fmt.Sprintf("luet tree pkglist -b -m \"%s\" --revdeps -o json --tree %s", revdep, dir))
	if err != nil {
		return []LuetPackage{}, errors.Wrap(err, "failed getting packages with luet")
	}

	packages := &TreeResult{}
	if err := json.Unmarshal([]byte(jsonPacks), packages); err != nil {
		return []LuetPackage{}, errors.Wrap(err, "failed getting packages with luet")
	}

	return packages.Packages, nil
}
