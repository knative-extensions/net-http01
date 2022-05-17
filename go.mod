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
	knative.dev/hack v0.0.0-20220512014059-f4972b4daff9
	knative.dev/networking v0.0.0-20220516202057-7363ea66ba54
	knative.dev/pkg v0.0.0-20220512013937-2d8305b2e59a
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
