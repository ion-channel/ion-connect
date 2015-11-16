package ionconnect

import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"
  "testing"
)

func TestIonConnect(t *testing.T) {
  RegisterFailHandler(Fail)
  RunSpecs(t, "Ion Connect Test Suite")
}
