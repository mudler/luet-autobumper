package autobumper_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAutobumper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Autobumper test suite")
}
