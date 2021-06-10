module knative.dev/net-http01

go 1.14

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/google/go-cmp v0.5.6
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	k8s.io/api v0.20.7
	k8s.io/apimachinery v0.20.7
	k8s.io/client-go v0.20.7
	knative.dev/hack v0.0.0-20210610231243-3d4b264d9472
	knative.dev/networking v0.0.0-20210610142944-8c7fb22941cf
	knative.dev/pkg v0.0.0-20210610171444-c96c7beb09df
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
