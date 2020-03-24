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
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
	"knative.dev/net-http01/pkg/challenger"
	"knative.dev/net-http01/pkg/ordermanager"
)

var domain = flag.String("domain", "", "The domain on which to serve the sample.")

func main() {
	flag.Parse()

	domains := []string{*domain}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Uncomment to use the Let's Encrypt staging endpoint.
	// ordermanager.Endpoint = ordermanager.Staging

	// Start our HTTP server to serve challenges.
	eg := errgroup.Group{}
	chlr, err := challenger.New(ctx)
	if err != nil {
		log.Fatalf("Error creating challenger: %v", err)
	}
	eg.Go(func() error { return http.ListenAndServe(":8080", chlr) })

	// Create our OrderManager, and provide a callback to signal us when
	// the certificate is ready to be picked up.  Give it our Challenger
	// to use for handling the HTTP01 challenges.
	ready := make(chan struct{})
	om, err := ordermanager.New(ctx, func(interface{}) {
		log.Print("Certificate should be ready!")
		close(ready)
	}, chlr)
	if err != nil {
		log.Fatalf("Error creating OrderManager: %v", err)
	}

	// First call returns the challenges (for us to set up Ingress)
	challs, _, err := om.Order(ctx, domains, nil)
	if err != nil {
		log.Fatalf("Error placing Domain order: %v", err)
	}
	log.Printf("Got challenges: %v", challs)

	// When the callback indicates that the certificate is ready, continue.
	select {
	case <-ready:
	case <-time.After(1 * time.Minute):
		log.Fatal("Timed out waiting for ready signal.")
	}

	// Calling order after the certificate is ready should yield the certificate.
	_, cert, err := om.Order(ctx, domains, nil)
	if err != nil {
		log.Fatalf("Error placing Domain order: %v", err)
	}

	s := &http.Server{
		Addr:      ":8443",
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{*cert}},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(w, "Hi")
		}),
	}

	eg.Go(func() error { return s.ListenAndServeTLS("", "") })

	if err := eg.Wait(); err != nil {
		log.Fatalf("Error serving: %v", err)
	}
}
