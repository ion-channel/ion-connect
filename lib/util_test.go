package ionconnect

import (
  . "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/onsi/ginkgo"
  . "github.com/ion-channel/ion-connect/Godeps/_workspace/src/github.com/onsi/gomega"
)

var _ = Describe("Util", func() {
  var (

  )

  BeforeEach(func() {
    Debug = true
  })

  Context("When the debug flag is set", func() {
    It("should write out debug statements", func() {
        Expect(func(){Debugln("testing")}).ShouldNot(Panic())
        Expect(func(){Debugf("testing %s", "f")}).ShouldNot(Panic())
    })
  })

  AfterEach(func() {
    Debug = false
  })
})
