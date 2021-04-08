module knative.dev/net-http01

go 1.14

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/google/go-cmp v0.5.5
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v0.19.7
	knative.dev/hack v0.0.0-20210325223819-b6ab329907d3
	knative.dev/networking v0.0.0-20210408132050-c8c1ee6e1873
	knative.dev/pkg v0.0.0-20210408073950-01dccc570bb3
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
