/*
Copyright 2020 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License"); you
may not use this file except in compliance with the License.  You may
obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
implied.  See the License for the specific language governing
permissions and limitations under the License.
*/

package resources

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/networking/pkg/apis/networking/v1alpha1"
	"knative.dev/pkg/ptr"
)

func TestMakeSecret(t *testing.T) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), cryptorand.Reader)
	if err != nil {
		t.Fatalf("ecdsa.GenerateKey() = %v", err)
	}
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := cryptorand.Int(cryptorand.Reader, serialNumberLimit)
	if err != nil {
		t.Fatalf("Failed to generate serial number: %v", err)
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Knative Ingress Conformance Testing"},
		},

		// Only let it live briefly.
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(5 * time.Minute),

		IsCA:                  true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"example.com"},
	}
	derBytes, err := x509.CreateCertificate(cryptorand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		t.Fatalf("x509.CreateCertificate() = %v", err)
	}
	der := [][]byte{derBytes}

	certPEM := &bytes.Buffer{}
	if err := pem.Encode(certPEM, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		t.Fatalf("pem.Encode() = %v", err)
	}
	result := certPEM.Bytes()

	tests := []struct {
		name string
		o    *v1alpha1.Certificate
		cert *tls.Certificate

		wantErr bool
		want    *corev1.Secret
	}{{
		name: "no private key",
		o: &v1alpha1.Certificate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
			Spec: v1alpha1.CertificateSpec{
				SecretName: "baz",
			},
		},
		cert: &tls.Certificate{
			Certificate: der,
			Leaf: &x509.Certificate{
				Raw: derBytes,
			},
			PrivateKey: nil,
		},
		wantErr: true,
	}, {
		name: "good secret",
		o: &v1alpha1.Certificate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
			Spec: v1alpha1.CertificateSpec{
				SecretName: "baz",
			},
		},
		cert: &tls.Certificate{
			Certificate: der,
			Leaf: &x509.Certificate{
				Raw: derBytes,
			},
			PrivateKey: priv,
		},
		want: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "baz",
				Namespace: "bar",
				OwnerReferences: []metav1.OwnerReference{{
					APIVersion:         "networking.internal.knative.dev/v1alpha1",
					Kind:               "Certificate",
					Name:               "foo",
					Controller:         ptr.Bool(true),
					BlockOwnerDeletion: ptr.Bool(true),
				}},
			},
			Type: corev1.SecretTypeTLS,
			Data: map[string][]byte{
				corev1.TLSCertKey:       result,
				corev1.TLSPrivateKeyKey: []byte(""),
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := MakeSecret(test.o, test.cert)
			switch {
			case test.wantErr:
				if err == nil {
					t.Errorf("MakeSecret() = %+v, wanted error", got)
				}
			case err != nil:
				t.Errorf("MakeSecret() = %v", err)
			default:
				if _, ok := got.Data[corev1.TLSPrivateKeyKey]; !ok {
					t.Errorf("Secret is missing key: %s", corev1.TLSPrivateKeyKey)
				}
				// TODO(mattmoor): Further validate the private key?

				// Clear it out, it's going to change every time.
				got.Data[corev1.TLSPrivateKeyKey] = []byte("")

				if !cmp.Equal(got, test.want) {
					t.Errorf("MakeSecret (-want, +got) = %s", cmp.Diff(got, test.want))
				}
			}
		})
	}
}

// Based on ./test/conformance/ingress/util.go#L700-L701 in knative/serving
func makeCert(t *testing.T, domains []string, expiry time.Time) []byte {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), cryptorand.Reader)
	if err != nil {
		t.Fatalf("ecdsa.GenerateKey() = %v", err)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := cryptorand.Int(cryptorand.Reader, serialNumberLimit)
	if err != nil {
		t.Fatalf("Failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Knative Ingress Conformance Testing"},
		},

		// Only let it live briefly.
		NotBefore: expiry.Add(-20 * time.Hour),
		NotAfter:  expiry,

		IsCA:                  true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,

		DNSNames: domains,
	}

	derBytes, err := x509.CreateCertificate(cryptorand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		t.Fatalf("x509.CreateCertificate() = %v", err)
	}

	certPEM := &bytes.Buffer{}
	if err := pem.Encode(certPEM, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		t.Fatalf("Failed to write data to cert.pem: %s", err)
	}

	return certPEM.Bytes()
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name            string
		secret          *corev1.Secret
		domains         []string
		minimumLifespan time.Duration
		want            *bool // nil means we want an error
	}{{
		name: "empty secret is invalid",
		secret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
		},
		want: ptr.Bool(false),
	}, {
		name: "no tls.crt key",
		secret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
			Data: map[string][]byte{
				// No tls.crt key
			},
		},
		want: ptr.Bool(false),
	}, {
		name: "bad PEM encoding",
		secret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
			Data: map[string][]byte{
				corev1.TLSCertKey: []byte("bad PEM data"),
			},
		},
		want: nil, // want an error
	}, {
		name: "good PEM, bad x509 cert",
		secret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
			Data: map[string][]byte{
				corev1.TLSCertKey: []byte("-----BEGIN CERTIFICATE-----\nZ2FyYmFnZQ==\n-----END CERTIFICATE-----\n"),
			},
		},
		want: nil, // want an error
	}, {
		name:            "missing domain",
		domains:         []string{"foo.com", "example.com"},
		minimumLifespan: 10 * time.Minute,
		secret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
			Data: map[string][]byte{
				corev1.TLSCertKey: makeCert(t, []string{"example.com"}, time.Now().Add(20*time.Minute)),
			},
		},
		want: ptr.Bool(false),
	}, {
		name:            "insufficient time left",
		domains:         []string{"example.com"},
		minimumLifespan: 10 * time.Minute,
		secret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
			Data: map[string][]byte{
				corev1.TLSCertKey: makeCert(t, []string{"example.com"}, time.Now().Add(5*time.Minute)),
			},
		},
		want: ptr.Bool(false),
	}, {
		name:            "good cert, single domain",
		domains:         []string{"example.com"},
		minimumLifespan: 10 * time.Minute,
		secret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
			Data: map[string][]byte{
				corev1.TLSCertKey: makeCert(t, []string{"example.com"}, time.Now().Add(30*time.Minute)),
			},
		},
		want: ptr.Bool(true),
	}, {
		name:            "good cert, multi-domain",
		domains:         []string{"example.com", "mattmoor.io"},
		minimumLifespan: 10 * time.Minute,
		secret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
			Data: map[string][]byte{
				corev1.TLSCertKey: makeCert(t, []string{"mattmoor.io", "example.com"}, time.Now().Add(30*time.Minute)),
			},
		},
		want: ptr.Bool(true),
	}, {
		name:            "good cert, extra domain",
		domains:         []string{"mattmoor.io"},
		minimumLifespan: 10 * time.Minute,
		secret: &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
			Data: map[string][]byte{
				corev1.TLSCertKey: makeCert(t, []string{"mattmoor.io", "example.com"}, time.Now().Add(30*time.Minute)),
			},
		},
		want: ptr.Bool(true),
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ok, err := IsValidCertificate(test.secret, test.domains, test.minimumLifespan)
			switch {
			case test.want == nil:
				if err == nil {
					t.Errorf("IsValidCertificate() = %v, wanted error", ok)
				}
			case err != nil:
				t.Errorf("IsValidCertificate() = %v", err)
			case *test.want != ok:
				t.Errorf("IsValidCertificate() = %v, wanted %v", ok, *test.want)
			}
		})
	}
}
