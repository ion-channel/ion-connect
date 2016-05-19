// ion-connect_suite_test.go
//
// Copyright (C) 2015 Selection Pressure LLC
//
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

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
