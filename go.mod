module knative.dev/net-http01

go 1.14

require (
	github.com/google/go-cmp v0.5.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.3 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	istio.io/gogo-genproto v0.0.0-20200130224810-a0338448499a // indirect
	k8s.io/api v0.18.1
	k8s.io/apimachinery v0.18.5
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/utils v0.0.0-20200320200009-4a6ff033650d // indirect
	knative.dev/networking v0.0.0-20200713162319-e2731eead7e8
	knative.dev/pkg v0.0.0-20200714070918-ac02cac99b88
	knative.dev/serving v0.16.1-0.20200714111218-995c90029adb
	knative.dev/test-infra v0.0.0-20200713220518-5a4c4cad5372
)

replace (
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
	k8s.io/api => k8s.io/api v0.17.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
	k8s.io/code-generator => k8s.io/code-generator v0.17.6
)
