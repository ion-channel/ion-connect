package leafnodes_test

import (
	. "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestLeafNode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LeafNode Suite")
}
