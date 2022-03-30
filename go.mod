module knative.dev/net-http01

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	github.com/imdario/mergo v0.3.9 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v0.23.5
	knative.dev/hack v0.0.0-20220328133751-f06773764ce3
	knative.dev/networking v0.0.0-20220329175515-09072d975ff9
	knative.dev/pkg v0.0.0-20220329144915-0a1ec2e0d46c
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
