/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package certificate

import (
	"testing"

	"knative.dev/net-http01/pkg/challenger"
	"knative.dev/net-http01/pkg/ordermanager"
	configmap "knative.dev/pkg/configmap"

	. "knative.dev/pkg/reconciler/testing"
)

func TestNewController(t *testing.T) {
	ctx, _ := SetupFakeContext(t)
	configMapWatcher := configmap.NewStaticWatcher()

	chlr, err := challenger.New(ctx)
	if err != nil {
		t.Fatalf("challenger.New() = %v", err)
	}

	ordermanager.Endpoint = ordermanager.Staging
	defer func() {
		ordermanager.Endpoint = ordermanager.Production
	}()

	c := NewController(ctx, configMapWatcher, chlr)
	if c == nil {
		t.Fatal("Expected NewController to return a non-nil value")
	}
}
