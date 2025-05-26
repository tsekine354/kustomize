// Copyright 2022 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package krusty_test

import (
	"testing"

	kusttest_test "sigs.k8s.io/kustomize/api/testutils/kusttest"
)

func TestStableListOrderAfterStrategicMergePatch(t *testing.T) {

	th := kusttest_test.MakeHarness(t)
	th.WriteF("resource.yaml", `
apiVersion: v1
kind: Pod
metadata:
  name: example
spec:
  containers:
  - env:
    - name: ALFA
      value: alfa
    - name: BRAVO
      value: bravo
    - name: CHARLIE
      value: "$(ALFA)"
    name: example
`)

	th.WriteF("patch.yaml", `
apiVersion: v1
kind: Pod
metadata:
  name: example
spec:
  containers:
  - env:
    - name: CHARLIE
      value: "$(BRAVO)"
    - name: NEW_ITEM
      value: new_item
    name: example
`)
	th.WriteK(".", `
resources:
  - resource.yaml
patches:
  - path: patch.yaml
`)

	m := th.Run(".", th.MakeDefaultOptions())
	th.AssertActualEqualsExpected(m, `
apiVersion: v1
kind: Pod
metadata:
  name: example
spec:
  containers:
  - env:
    - name: NEW_ITEM
      value: new_item
    - name: ALFA
      value: alfa
    - name: BRAVO
      value: bravo
    - name: CHARLIE
      value: "$(BRAVO)"
    name: example
`)
}
