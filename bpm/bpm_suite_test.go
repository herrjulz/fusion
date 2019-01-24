package bpm_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestBpm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bpm Suite")
}
