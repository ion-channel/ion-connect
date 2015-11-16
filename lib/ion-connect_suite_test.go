package ionconnect

import (
  . "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/onsi/ginkgo"
  . "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/onsi/gomega"
  "testing"
)

func TestIonConnect(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "Ion Connect Test Suite")
}
