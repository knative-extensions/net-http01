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
	knative.dev/hack v0.0.0-20210602212444-509255f29a24
	knative.dev/networking v0.0.0-20210608114541-4b1712c029b7
	knative.dev/pkg v0.0.0-20210902175106-8d4b5e065ebb
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
