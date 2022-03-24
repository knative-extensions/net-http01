module knative.dev/net-http01

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	github.com/imdario/mergo v0.3.9 // indirect
	golang.org/x/crypto v0.0.0-20220214200702-86341886e292
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	k8s.io/api v0.23.4
	k8s.io/apimachinery v0.23.4
	k8s.io/client-go v0.23.4
	knative.dev/hack v0.0.0-20220318020218-14f832e506f8
	knative.dev/networking v0.0.0-20220323170318-55757e9c20d6
	knative.dev/pkg v0.0.0-20220318185521-e6e3cf03d765
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
