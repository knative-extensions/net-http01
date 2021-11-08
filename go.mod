module knative.dev/net-http01

go 1.14

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/google/go-cmp v0.5.6
	golang.org/x/crypto v0.0.0-20210920023735-84f357641f63
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	k8s.io/api v0.21.4
	k8s.io/apimachinery v0.21.4
	k8s.io/client-go v0.21.4
	knative.dev/hack v0.0.0-20211105231158-29f86c2653b5
	knative.dev/networking v0.0.0-20211104064801-6871f98f7b4d
	knative.dev/pkg v0.0.0-20211104101302-51b9e7f161b4
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
