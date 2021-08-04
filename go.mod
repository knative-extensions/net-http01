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
	knative.dev/hack v0.0.0-20210622141627-e28525d8d260
	knative.dev/networking v0.0.0-20210803181815-acdfd41c575c
	knative.dev/pkg v0.0.0-20210803160015-21eb4c167cc5
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
