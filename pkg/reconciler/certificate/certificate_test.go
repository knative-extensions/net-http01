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
	"bytes"
	context "context"
	"crypto/ecdsa"
	"crypto/elliptic"
	cryptorand "crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	clientgotesting "k8s.io/client-go/testing"
	"knative.dev/net-http01/pkg/ordermanager"
	"knative.dev/net-http01/pkg/reconciler/certificate/resources"
	"knative.dev/networking/pkg/apis/networking"
	v1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	certreconciler "knative.dev/networking/pkg/client/injection/reconciler/networking/v1alpha1/certificate"
	"knative.dev/pkg/apis"
	configmap "knative.dev/pkg/configmap"
	controller "knative.dev/pkg/controller"
	logging "knative.dev/pkg/logging"

	networkingclient "knative.dev/networking/pkg/client/injection/client/fake"
	_ "knative.dev/networking/pkg/client/injection/informers/networking/v1alpha1/certificate/fake"
	kubeclient "knative.dev/pkg/client/injection/kube/client/fake"
	_ "knative.dev/pkg/client/injection/kube/informers/core/v1/endpoints/fake"
	_ "knative.dev/pkg/client/injection/kube/informers/core/v1/secret/fake"
	_ "knative.dev/pkg/client/injection/kube/informers/core/v1/service/fake"

	. "knative.dev/net-http01/pkg/reconciler/testing"
	. "knative.dev/pkg/reconciler/testing"
)

func TestReconcileMakingOrders(t *testing.T) {
	table := TableTest{{
		Name: "bad workqueue key",
		Key:  "too/many/parts",
	}, {
		Name: "create a bunch of stuff, make order, set challenges",
		Objects: []runtime.Object{
			cert("kn-cert", "foo", withDomains("example.com")),
		},
		WantCreates: []runtime.Object{
			resources.MakeService(cert("kn-cert", "foo", withDomains("example.com"))),
			resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com"))),
		},
		WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
			Object: cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.MarkNotReady("OrderCert", "Provisioning Certificate through HTTP01 challenges.")
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
				}),
		}},
		Key: "foo/kn-cert",
	}, {
		Name: "steady state post creation",
		Objects: []runtime.Object{
			cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.MarkNotReady("OrderCert", "Provisioning Certificate through HTTP01 challenges.")
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
				}),
			resources.MakeService(cert("kn-cert", "foo", withDomains("example.com"))),
			resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com"))),
		},
		Key: "foo/kn-cert",
	}, {
		Name: "cert has dots",
		Objects: []runtime.Object{
			cert("kn.cert.io", "foo", withDomains("example.com"), withUID("42-42-42"),
				func(c *v1alpha1.Certificate) {
					c.Status.MarkNotReady("OrderCert", "Provisioning Certificate through HTTP01 challenges.")
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
				}),
			resources.MakeService(cert("kn.cert.io", "foo", withDomains("example.com"), withUID("42-42-42"))),
			resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com"), withUID("42-42-42"))),
		},
		Key: "foo/kn-cert",
	}, {
		Name: "update bad service",
		Objects: []runtime.Object{
			cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.MarkNotReady("OrderCert", "Provisioning Certificate through HTTP01 challenges.")
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
				}),
			resources.MakeService(cert("kn-cert", "foo", withDomains("example.com")),
				func(svc *corev1.Service) {
					svc.Spec = corev1.ServiceSpec{}
				}),
			resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com"))),
		},
		WantUpdates: []clientgotesting.UpdateActionImpl{{
			Object: resources.MakeService(cert("kn-cert", "foo", withDomains("example.com"))),
		}},
		Key: "foo/kn-cert",
	}, {
		Name: "update bad endpoints",
		Objects: []runtime.Object{
			cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.MarkNotReady("OrderCert", "Provisioning Certificate through HTTP01 challenges.")
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
				}),
			resources.MakeService(cert("kn-cert", "foo", withDomains("example.com"))),
			resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com")),
				func(ep *corev1.Endpoints) {
					ep.Subsets = nil
				}),
		},
		WantUpdates: []clientgotesting.UpdateActionImpl{{
			Object: resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com"))),
		}},
		Key: "foo/kn-cert",
	}, {
		Name:    "error creating service",
		WantErr: true,
		WithReactors: []clientgotesting.ReactionFunc{
			InduceFailure("create", "services"),
		},
		Objects: []runtime.Object{
			cert("kn-cert", "foo", withDomains("example.com")),
		},
		WantCreates: []runtime.Object{
			resources.MakeService(cert("kn-cert", "foo", withDomains("example.com"))),
		},
		WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
			Object: cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.InitializeConditions()
				}),
		}},
		WantEvents: []string{
			Eventf(corev1.EventTypeWarning, "InternalError", "inducing failure for create services"),
		},
		Key: "foo/kn-cert",
	}, {
		Name:    "error creating endpoints",
		WantErr: true,
		WithReactors: []clientgotesting.ReactionFunc{
			InduceFailure("create", "endpoints"),
		},
		Objects: []runtime.Object{
			cert("kn-cert", "foo", withDomains("example.com")),
		},
		WantCreates: []runtime.Object{
			resources.MakeService(cert("kn-cert", "foo", withDomains("example.com"))),
			resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com"))),
		},
		WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
			Object: cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.InitializeConditions()
				}),
		}},
		WantEvents: []string{
			Eventf(corev1.EventTypeWarning, "InternalError", "inducing failure for create endpoints"),
		},
		Key: "foo/kn-cert",
	}, {
		Name:    "error updating service",
		WantErr: true,
		WithReactors: []clientgotesting.ReactionFunc{
			InduceFailure("update", "services"),
		},
		Objects: []runtime.Object{
			cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.MarkNotReady("OrderCert", "Provisioning Certificate through HTTP01 challenges.")
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
				}),
			resources.MakeService(cert("kn-cert", "foo", withDomains("example.com")),
				func(svc *corev1.Service) {
					svc.Spec = corev1.ServiceSpec{}
				}),
			resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com"))),
		},
		WantUpdates: []clientgotesting.UpdateActionImpl{{
			Object: resources.MakeService(cert("kn-cert", "foo", withDomains("example.com"))),
		}},
		WantEvents: []string{
			Eventf(corev1.EventTypeWarning, "InternalError", "inducing failure for update services"),
		},
		Key: "foo/kn-cert",
	}, {
		Name:    "error updating endpoints",
		WantErr: true,
		WithReactors: []clientgotesting.ReactionFunc{
			InduceFailure("update", "endpoints"),
		},
		Objects: []runtime.Object{
			cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.MarkNotReady("OrderCert", "Provisioning Certificate through HTTP01 challenges.")
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
				}),
			resources.MakeService(cert("kn-cert", "foo", withDomains("example.com"))),
			resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com")),
				func(ep *corev1.Endpoints) {
					ep.Subsets = nil
				}),
		},
		WantUpdates: []clientgotesting.UpdateActionImpl{{
			Object: resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com"))),
		}},
		WantEvents: []string{
			Eventf(corev1.EventTypeWarning, "InternalError", "inducing failure for update endpoints"),
		},
		Key: "foo/kn-cert",
	}, {
		Name: "valid secret",
		Objects: []runtime.Object{
			cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.MarkNotReady("OrderCert", "Provisioning Certificate through HTTP01 challenges.")
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
				}),
			resources.MakeService(cert("kn-cert", "foo", withDomains("example.com"))),
			resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com"))),
			// This is based on the IsValidCert unit test.
			&corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kn-cert",
					Namespace: "foo",
				},
				Data: map[string][]byte{
					corev1.TLSCertKey: makeCert(t, []string{"example.com"}, time.Now().Add(100*24*time.Hour)),
				},
			},
		},
		WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
			Object: cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.InitializeConditions()
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
					// Becomes ready.
					c.Status.MarkReady()
				}),
		}},
		Key: "foo/kn-cert",
	}, {
		Name: "not enough time left",
		Objects: []runtime.Object{
			cert("kn-cert", "foo", withDomains("example.com")),
			resources.MakeService(cert("kn-cert", "foo", withDomains("example.com"))),
			resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com"))),
			// This is based on the IsValidCert unit test.
			&corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kn-cert",
					Namespace: "foo",
				},
				Data: map[string][]byte{
					corev1.TLSCertKey: makeCert(t, []string{"example.com"}, time.Now().Add(1*time.Hour)),
				},
			},
		},
		WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
			Object: cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.InitializeConditions()
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
					c.Status.MarkNotReady("OrderCert", "Provisioning Certificate through HTTP01 challenges.")
				}),
		}},
		Key: "foo/kn-cert",
	}}

	table.Test(t, MakeFactory(func(ctx context.Context, listers *Listers, cmw configmap.Watcher) controller.Reconciler {
		r := &Reconciler{
			kubeClient:      kubeclient.Get(ctx),
			secretLister:    listers.GetSecretLister(),
			serviceLister:   listers.GetK8sServiceLister(),
			endpointsLister: listers.GetEndpointsLister(),
			challengePort:   8080,

			orderManager: &fakeOM{
				challenges: []*apis.URL{{
					Scheme: "http",
					Host:   "example.com",
					Path:   "/.acme/well-known/gobbledy-gook",
				}},
			},
		}

		return certreconciler.NewReconciler(ctx, logging.FromContext(ctx), networkingclient.Get(ctx),
			listers.GetCertificateLister(), controller.GetEventRecorder(ctx), r, CertificateClassName)
	}))
}

func TestReconcileOrderError(t *testing.T) {
	table := TableTest{{
		Name:    "error placing order",
		WantErr: true,
		Objects: []runtime.Object{
			cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.MarkNotReady("OrderCert", "Provisioning Certificate through HTTP01 challenges.")
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
				}),
			resources.MakeService(cert("kn-cert", "foo", withDomains("example.com"))),
			resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com"))),
		},
		Key: "foo/kn-cert",
		WantEvents: []string{
			Eventf(corev1.EventTypeWarning, "InternalError", "an error"),
		},
	}}

	table.Test(t, MakeFactory(func(ctx context.Context, listers *Listers, cmw configmap.Watcher) controller.Reconciler {
		r := &Reconciler{
			kubeClient:      kubeclient.Get(ctx),
			secretLister:    listers.GetSecretLister(),
			serviceLister:   listers.GetK8sServiceLister(),
			endpointsLister: listers.GetEndpointsLister(),
			challengePort:   8080,

			orderManager: &fakeOM{
				err: errors.New("an error"),
			},
		}

		return certreconciler.NewReconciler(ctx, logging.FromContext(ctx), networkingclient.Get(ctx),
			listers.GetCertificateLister(), controller.GetEventRecorder(ctx), r, CertificateClassName)
	}))
}

func TestReconcileOrderFulfillment(t *testing.T) {

	tc := makeTLSCert(t, []string{"example.com"}, time.Now().Add(100*24*time.Hour))

	table := TableTest{{
		Name: "create a new TLS secret",
		Objects: []runtime.Object{
			cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.MarkNotReady("OrderCert", "Provisioning Certificate through HTTP01 challenges.")
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
				}),
			resources.MakeService(cert("kn-cert", "foo", withDomains("example.com"))),
			resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com"))),
		},
		WantCreates: []runtime.Object{
			mustMakeSecret(t, cert("kn-cert", "foo"), tc),
		},
		WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
			Object: cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.InitializeConditions()
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
					c.Status.MarkReady()
				}),
		}},
		Key: "foo/kn-cert",
	}, {
		Name: "update the TLS secret",
		Objects: []runtime.Object{
			cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.MarkNotReady("OrderCert", "Provisioning Certificate through HTTP01 challenges.")
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
				}),
			resources.MakeService(cert("kn-cert", "foo", withDomains("example.com"))),
			resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com"))),
			mustMakeSecret(t, cert("kn-cert", "foo"), tc, func(s *corev1.Secret) {
				s.Data = nil
			}),
		},
		WantUpdates: []clientgotesting.UpdateActionImpl{{
			Object: mustMakeSecret(t, cert("kn-cert", "foo"), tc),
		}},
		WantStatusUpdates: []clientgotesting.UpdateActionImpl{{
			Object: cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.InitializeConditions()
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
					c.Status.MarkReady()
				}),
		}},
		Key: "foo/kn-cert",
	}, {
		Name:    "error creating secret",
		WantErr: true,
		WithReactors: []clientgotesting.ReactionFunc{
			InduceFailure("create", "secrets"),
		},
		Objects: []runtime.Object{
			cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.MarkNotReady("OrderCert", "Provisioning Certificate through HTTP01 challenges.")
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
				}),
			resources.MakeService(cert("kn-cert", "foo", withDomains("example.com"))),
			resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com"))),
		},
		WantCreates: []runtime.Object{
			mustMakeSecret(t, cert("kn-cert", "foo"), tc),
		},
		WantEvents: []string{
			Eventf(corev1.EventTypeWarning, "InternalError", "inducing failure for create secrets"),
		},
		Key: "foo/kn-cert",
	}, {
		Name:    "error updating secret",
		WantErr: true,
		WithReactors: []clientgotesting.ReactionFunc{
			InduceFailure("update", "secrets"),
		},
		Objects: []runtime.Object{
			cert("kn-cert", "foo", withDomains("example.com"),
				func(c *v1alpha1.Certificate) {
					c.Status.MarkNotReady("OrderCert", "Provisioning Certificate through HTTP01 challenges.")
					c.Status.HTTP01Challenges = []v1alpha1.HTTP01Challenge{{
						ServiceName:      "kn-cert",
						ServiceNamespace: "foo",
						ServicePort:      intstr.FromInt(80),
						URL: &apis.URL{
							Scheme: "http",
							Host:   "example.com",
							Path:   "/.acme/well-known/gobbledy-gook",
						},
					}}
				}),
			resources.MakeService(cert("kn-cert", "foo", withDomains("example.com"))),
			resources.MakeEndpoints(cert("kn-cert", "foo", withDomains("example.com"))),
			mustMakeSecret(t, cert("kn-cert", "foo"), tc, func(s *corev1.Secret) {
				s.Data = nil
			}),
		},
		WantUpdates: []clientgotesting.UpdateActionImpl{{
			Object: mustMakeSecret(t, cert("kn-cert", "foo"), tc),
		}},
		WantEvents: []string{
			Eventf(corev1.EventTypeWarning, "InternalError", "inducing failure for update secrets"),
		},
		Key: "foo/kn-cert",
	}}

	table.Test(t, MakeFactory(func(ctx context.Context, listers *Listers, cmw configmap.Watcher) controller.Reconciler {
		r := &Reconciler{
			kubeClient:      kubeclient.Get(ctx),
			secretLister:    listers.GetSecretLister(),
			serviceLister:   listers.GetK8sServiceLister(),
			endpointsLister: listers.GetEndpointsLister(),
			challengePort:   8080,

			orderManager: &fakeOM{
				cert: tc,
			},
		}

		return certreconciler.NewReconciler(ctx, logging.FromContext(ctx), networkingclient.Get(ctx),
			listers.GetCertificateLister(), controller.GetEventRecorder(ctx), r, CertificateClassName)
	}))
}

type certOption func(*v1alpha1.Certificate)

func cert(name, namespace string, opts ...certOption) *v1alpha1.Certificate {
	c := &v1alpha1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Annotations: map[string]string{
				networking.CertificateClassAnnotationKey: CertificateClassName,
			},
		},
		Spec: v1alpha1.CertificateSpec{
			SecretName: name,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func withDomains(domains ...string) certOption {
	return func(c *v1alpha1.Certificate) {
		c.Spec.DNSNames = domains
	}
}

func withUID(uid types.UID) certOption {
	return func(c *v1alpha1.Certificate) {
		c.UID = uid
	}
}

type fakeOM struct {
	challenges []*apis.URL
	cert       *tls.Certificate
	err        error
}

var _ ordermanager.Interface = (*fakeOM)(nil)

func (fom *fakeOM) Order(ctx context.Context, domains []string, owner interface{}) ([]*apis.URL, *tls.Certificate, error) {
	switch {
	case fom.challenges != nil:
		return fom.challenges, nil, nil
	case fom.cert != nil:
		return nil, fom.cert, nil
	case fom.err != nil:
		return nil, nil, fom.err
	default:
		panic("fakeOM was improperly configured")
	}
}

func mustMakeSecret(t *testing.T, o *v1alpha1.Certificate, cert *tls.Certificate, opts ...func(*corev1.Secret)) *corev1.Secret {
	s, err := resources.MakeSecret(o, cert)
	if err != nil {
		t.Fatalf("MakeSecret() = %v", err)
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func makeTLSCert(t *testing.T, domains []string, expiry time.Time) *tls.Certificate {
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

	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		t.Fatalf("x509.ParseCertificate() = %v", err)
	}

	return &tls.Certificate{
		Certificate: [][]byte{derBytes},
		Leaf:        cert,
		PrivateKey:  priv,
	}
}

// Based on ./test/conformance/ingress/util.go#L700-L701 in knative/serving
func makeCert(t *testing.T, domains []string, expiry time.Time) []byte {
	tc := makeTLSCert(t, domains, expiry)

	certPEM := &bytes.Buffer{}
	if err := pem.Encode(certPEM, &pem.Block{Type: "CERTIFICATE", Bytes: tc.Leaf.Raw}); err != nil {
		t.Fatalf("Failed to write data to cert.pem: %s", err)
	}

	return certPEM.Bytes()
}
