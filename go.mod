module knative.dev/net-http01

go 1.14

require (
	github.com/google/go-cmp v0.5.1
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	k8s.io/api v0.18.7-rc.0
	k8s.io/apimachinery v0.18.7-rc.0
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/networking v0.0.0-20200812200006-4d518e76538a
	knative.dev/pkg v0.0.0-20200812224206-44c860147a87
	knative.dev/test-infra v0.0.0-20200813220834-388e55a496cf
)

replace (
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
	k8s.io/api => k8s.io/api v0.17.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
	k8s.io/code-generator => k8s.io/code-generator v0.17.6
)
