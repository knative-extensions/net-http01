module knative.dev/net-http01

go 1.14

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/google/go-cmp v0.5.6
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	knative.dev/hack v0.0.0-20211018110626-47ac3b032e60
	knative.dev/networking v0.0.0-20211018053525-91e7d042db3d
	knative.dev/pkg v0.0.0-20211018071527-3463dd75659d
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
