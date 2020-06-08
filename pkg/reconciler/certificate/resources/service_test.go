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
	"testing"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/networking/pkg/apis/networking/v1alpha1"
	"knative.dev/pkg/ptr"
)

func TestMakeService(t *testing.T) {
	tests := []struct {
		name string
		o    *v1alpha1.Certificate
		want *corev1.Service
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
			Spec: serviceSpec,
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
			Spec: serviceSpec,
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := MakeService(test.o)
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
			Subsets: endpointSubsets,
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
			Subsets: endpointSubsets,
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := MakeEndpoints(test.o)
			if !cmp.Equal(got, test.want) {
				t.Errorf("MakeEndpoints (-want, +got) = %s", cmp.Diff(got, test.want))
			}
		})
	}
}
