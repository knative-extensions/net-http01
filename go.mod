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
	knative.dev/hack v0.0.0-20210114150620-4422dcadb3c8
	knative.dev/networking v0.0.0-20210118031120-e223a8941d09
	knative.dev/pkg v0.0.0-20210118142921-22f4898e6ef7
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
