package mytheresa_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMytheresa(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mytheresa Suite")
}
