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
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"knative.dev/networking/pkg/apis/networking/v1alpha1"
	"knative.dev/pkg/ptr"
)

func TestMakeService(t *testing.T) {
	tests := []struct {
		name string
		o    *v1alpha1.Certificate
		want *corev1.Service
		opts []func(*corev1.Service)
	}{{
		name: "check owner refs one way",
		o: &v1alpha1.Certificate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
		},
		want: &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
				OwnerReferences: []metav1.OwnerReference{{
					APIVersion:         "networking.internal.knative.dev/v1alpha1",
					Kind:               "Certificate",
					Name:               "foo",
					Controller:         ptr.Bool(true),
					BlockOwnerDeletion: ptr.Bool(true),
				}},
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{{
					Name:       portName,
					Port:       80,
					TargetPort: intstr.FromInt(8080),
				}},
			},
		},
	}, {
		name: "check owner refs another way",
		o: &v1alpha1.Certificate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar",
				Namespace: "food",
			},
		},
		want: &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar",
				Namespace: "food",
				OwnerReferences: []metav1.OwnerReference{{
					APIVersion:         "networking.internal.knative.dev/v1alpha1",
					Kind:               "Certificate",
					Name:               "bar",
					Controller:         ptr.Bool(true),
					BlockOwnerDeletion: ptr.Bool(true),
				}},
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{{
					Name:       portName,
					Port:       80,
					TargetPort: intstr.FromInt(8080),
				}},
			},
		},
	}, {
		name: "custom port",
		o: &v1alpha1.Certificate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
		},
		want: &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
				OwnerReferences: []metav1.OwnerReference{{
					APIVersion:         "networking.internal.knative.dev/v1alpha1",
					Kind:               "Certificate",
					Name:               "foo",
					Controller:         ptr.Bool(true),
					BlockOwnerDeletion: ptr.Bool(true),
				}},
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{{
					Name:       portName,
					Port:       80,
					TargetPort: intstr.FromInt(1234),
				}},
			},
		},
		opts: []func(*corev1.Service){WithServicePort(1234)},
	}, {
		name: "name has dots",
		o: &v1alpha1.Certificate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo.com",
				Namespace: "bar",
				// UUID should be 35 chars long
				UID: "123e4567-e89b-12d3-a456-426614174000",
			},
		},
		want: &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "challenge-for-123e4567-e89b-12d3-a456-426614174000",
				Namespace: "bar",
				OwnerReferences: []metav1.OwnerReference{{
					APIVersion:         "networking.internal.knative.dev/v1alpha1",
					Kind:               "Certificate",
					Name:               "foo.com",
					UID:                "123e4567-e89b-12d3-a456-426614174000",
					Controller:         ptr.Bool(true),
					BlockOwnerDeletion: ptr.Bool(true),
				}},
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{{
					Name:       portName,
					Port:       80,
					TargetPort: intstr.FromInt(8080),
				}},
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := MakeService(test.o, test.opts...)
			if !cmp.Equal(got, test.want) {
				t.Errorf("MakeService (-want, +got) = %s", cmp.Diff(got, test.want))
			}
		})
	}
}

func TestMakeEndpoints(t *testing.T) {
	tests := []struct {
		name string
		o    *v1alpha1.Certificate
		want *corev1.Endpoints
		opts []func(*corev1.Endpoints)
	}{{
		name: "check owner refs one way",
		o: &v1alpha1.Certificate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
		},
		want: &corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
				OwnerReferences: []metav1.OwnerReference{{
					APIVersion:         "networking.internal.knative.dev/v1alpha1",
					Kind:               "Certificate",
					Name:               "foo",
					Controller:         ptr.Bool(true),
					BlockOwnerDeletion: ptr.Bool(true),
				}},
			},
			Subsets: []corev1.EndpointSubset{{
				Addresses: []corev1.EndpointAddress{{
					IP: os.Getenv("POD_IP"),
				}},
				Ports: []corev1.EndpointPort{{
					Name:     portName,
					Port:     8080,
					Protocol: corev1.ProtocolTCP,
				}},
			}},
		},
	}, {
		name: "check owner refs another way",
		o: &v1alpha1.Certificate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar",
				Namespace: "food",
			},
		},
		want: &corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar",
				Namespace: "food",
				OwnerReferences: []metav1.OwnerReference{{
					APIVersion:         "networking.internal.knative.dev/v1alpha1",
					Kind:               "Certificate",
					Name:               "bar",
					Controller:         ptr.Bool(true),
					BlockOwnerDeletion: ptr.Bool(true),
				}},
			},
			Subsets: []corev1.EndpointSubset{{
				Addresses: []corev1.EndpointAddress{{
					IP: os.Getenv("POD_IP"),
				}},
				Ports: []corev1.EndpointPort{{
					Name:     portName,
					Port:     8080,
					Protocol: corev1.ProtocolTCP,
				}},
			}},
		},
	}, {
		name: "custom port",
		o: &v1alpha1.Certificate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
			},
		},
		want: &corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "foo",
				Namespace: "bar",
				OwnerReferences: []metav1.OwnerReference{{
					APIVersion:         "networking.internal.knative.dev/v1alpha1",
					Kind:               "Certificate",
					Name:               "foo",
					Controller:         ptr.Bool(true),
					BlockOwnerDeletion: ptr.Bool(true),
				}},
			},
			Subsets: []corev1.EndpointSubset{{
				Addresses: []corev1.EndpointAddress{{
					IP: os.Getenv("POD_IP"),
				}},
				Ports: []corev1.EndpointPort{{
					Name:     portName,
					Port:     1234,
					Protocol: corev1.ProtocolTCP,
				}},
			}},
		},
		opts: []func(*corev1.Endpoints){WithEndpointsPort(1234)},
	}, {
		name: "name has dots",
		o: &v1alpha1.Certificate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "bar.com",
				Namespace: "food",
				UID:       "dead-beef",
			},
		},
		want: &corev1.Endpoints{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "challenge-for-dead-beef",
				Namespace: "food",
				OwnerReferences: []metav1.OwnerReference{{
					APIVersion:         "networking.internal.knative.dev/v1alpha1",
					Kind:               "Certificate",
					Name:               "bar.com",
					UID:                "dead-beef",
					Controller:         ptr.Bool(true),
					BlockOwnerDeletion: ptr.Bool(true),
				}},
			},
			Subsets: []corev1.EndpointSubset{{
				Addresses: []corev1.EndpointAddress{{
					IP: os.Getenv("POD_IP"),
				}},
				Ports: []corev1.EndpointPort{{
					Name:     portName,
					Port:     8080,
					Protocol: corev1.ProtocolTCP,
				}},
			}},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := MakeEndpoints(test.o, test.opts...)
			if !cmp.Equal(got, test.want) {
				t.Errorf("MakeEndpoints (-want, +got) = %s", cmp.Diff(got, test.want))
			}
		})
	}
}
