package utils

import (
	"github.com/mikefarah/yq/v4/pkg/yqlib"
	"gopkg.in/op/go-logging.v1"
)

func YQ(expr, file string) error {
	logging.SetLevel(logging.CRITICAL, "yq-lib")

	var completedSuccessfully bool
	format, err := yqlib.OutputFormatFromString("yaml")
	if err != nil {
		return err
	}
	writeInPlaceHandler := yqlib.NewWriteInPlaceHandler(file)
	out, err := writeInPlaceHandler.CreateTempFile()
	if err != nil {
		return err
	}
	// need to indirectly call the function so  that completedSuccessfully is
	// passed when we finish execution as opposed to now
	defer func() { writeInPlaceHandler.FinishWriteInPlace(completedSuccessfully) }()

	printer := yqlib.NewPrinter(out, format, false, false, 0, false)

	streamEvaluator := yqlib.NewStreamEvaluator()

	err = streamEvaluator.EvaluateFiles(expr, []string{file}, printer, true)
	completedSuccessfully = err == nil
	return err
}
