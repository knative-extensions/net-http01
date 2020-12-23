module knative.dev/net-http01

go 1.14

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/google/go-cmp v0.5.4
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/prometheus/procfs v0.0.11 // indirect
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	k8s.io/api v0.18.12
	k8s.io/apimachinery v0.18.12
	k8s.io/client-go v0.18.12
	knative.dev/hack v0.0.0-20201214230143-4ed1ecb8db24
	knative.dev/networking v0.0.0-20201223045804-ccb4f21cdf5c
	knative.dev/pkg v0.0.0-20201223002104-9d0775512af8
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
