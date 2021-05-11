module knative.dev/net-http01

go 1.14

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/google/go-cmp v0.5.5
	golang.org/x/crypto v0.0.0-20210415154028-4f45737414dc
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v0.19.7
	knative.dev/hack v0.0.0-20210428122153-93ad9129c268
	knative.dev/networking v0.0.0-20210506040209-a3028d57082a
	knative.dev/pkg v0.0.0-20210510175900-4564797bf3b7
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
