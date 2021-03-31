module knative.dev/net-http01

go 1.14

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/google/go-cmp v0.5.5
	golang.org/x/crypto v0.0.0-20201002170205-7f63de1d35b0
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	k8s.io/api v0.19.7
	k8s.io/apimachinery v0.19.7
	k8s.io/client-go v0.19.7
	knative.dev/hack v0.0.0-20210325223819-b6ab329907d3
	knative.dev/networking v0.0.0-20210331064822-999a7708876c
	knative.dev/pkg v0.0.0-20210331065221-952fdd90dbb0
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
