package containernode_test

import (
	. "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestContainernode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Containernode Suite")
}
