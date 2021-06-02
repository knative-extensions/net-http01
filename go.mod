module knative.dev/net-http01

go 1.14

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/google/go-cmp v0.5.6
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v0.19.7
	knative.dev/hack v0.0.0-20210601210329-de04b70e00d0
	knative.dev/networking v0.0.0-20210602143631-9c0fc00ae8fe
	knative.dev/pkg v0.0.0-20210602095030-0e61d6763dd6
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
