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

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	network "knative.dev/networking/pkg"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/signals"

	"knative.dev/net-http01/pkg/challenger"
	"knative.dev/net-http01/pkg/reconciler/certificate"
)

func main() {
	// Uncomment this to use the Let's Encrypt Staging environment.
	// ordermanager.Endpoint = ordermanager.Staging

	ctx := signals.NewContext()

	chlr, err := challenger.New(ctx)
	if err != nil {
		log.Fatalf("Error creating challenger: %v", err)
	}

	port := 8765

	go http.ListenAndServe(fmt.Sprint(":", port), network.NewProbeHandler(chlr))

	sharedmain.MainWithContext(ctx, "net-http01",
		func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
			return certificate.NewController(ctx, cmw, chlr, port)
		},
	)
}
