module knative.dev/net-http01

go 1.14

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/google/go-cmp v0.5.2
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/prometheus/procfs v0.0.11 // indirect
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	k8s.io/api v0.18.12
	k8s.io/apimachinery v0.18.12
	k8s.io/client-go v0.18.12
	knative.dev/hack v0.0.0-20201120192952-353db687ec5b
	knative.dev/networking v0.0.0-20201123014253-96ce58eb8018
	knative.dev/pkg v0.0.0-20201123014053-92bc25a0a520
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
