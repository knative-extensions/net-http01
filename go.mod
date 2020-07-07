module knative.dev/net-http01

go 1.14

require (
	github.com/google/go-cmp v0.5.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.3 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	istio.io/gogo-genproto v0.0.0-20200130224810-a0338448499a // indirect
	k8s.io/api v0.18.1
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/utils v0.0.0-20200320200009-4a6ff033650d // indirect
	knative.dev/networking v0.0.0-20200707203944-725ec013d8a2
	knative.dev/pkg v0.0.0-20200707190344-0a8314b44495
	knative.dev/serving v0.16.1-0.20200707204544-f025d519072b
	knative.dev/test-infra v0.0.0-20200707220947-a05ea93e9076
)

replace (
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
	k8s.io/api => k8s.io/api v0.17.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
	k8s.io/code-generator => k8s.io/code-generator v0.17.6
)
