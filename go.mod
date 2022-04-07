module knative.dev/net-http01

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	github.com/imdario/mergo v0.3.9 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	k8s.io/api v0.23.5
	k8s.io/apimachinery v0.23.5
	k8s.io/client-go v0.23.5
	knative.dev/hack v0.0.0-20220401031746-a75ca495e7f4
	knative.dev/networking v0.0.0-20220407031944-7fa8012b6f2d
	knative.dev/pkg v0.0.0-20220407131646-8db11d0f9fed
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
