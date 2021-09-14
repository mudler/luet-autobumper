package plugins

import (
	"fmt"

	"github.com/Luet-lab/luet-autobumper/pkg/autobumper"
	"github.com/mikefarah/yq/v4/pkg/yqlib"
)

type Inplace struct {
}

func (inplace *Inplace) Bump(src autobumper.LuetPackageWithLabels, dst autobumper.LuetPackageWithLabels) error {
	var completedSuccessfully bool
	var expr string
	if !src.IsCollection() {

		expr = fmt.Sprintf(".version = \"%s\"", dst.Version)
	} else {
		// 1: find the package inside the collection (index)
		coll, err := autobumper.ReadCollection(src.Path)
		if err != nil {
			return err
		}

		index, err := autobumper.Collection(coll).Find(src)
		if err != nil {
			return err
		}
		// 2: find the respective yaml.Node
		expr = fmt.Sprintf(".packages[%d].version = \"%s\"", index, dst.Version)

	}
	format, err := yqlib.OutputFormatFromString("yaml")
	if err != nil {
		return err
	}
	writeInPlaceHandler := yqlib.NewWriteInPlaceHandler(src.GetPath())
	out, err := writeInPlaceHandler.CreateTempFile()
	if err != nil {
		return err
	}
	// need to indirectly call the function so  that completedSuccessfully is
	// passed when we finish execution as opposed to now
	defer func() { writeInPlaceHandler.FinishWriteInPlace(completedSuccessfully) }()

	printer := yqlib.NewPrinter(out, format, false, false, 0, false)

	streamEvaluator := yqlib.NewStreamEvaluator()

	err = streamEvaluator.EvaluateFiles(expr, []string{src.GetPath()}, printer, true)
	completedSuccessfully = err == nil
	return err
}

func (inplace *Inplace) Apply(p autobumper.LuetPackageWithLabels) bool {
	return p.Labels["autobump.inplace"] != "false"
}
