module knative.dev/net-http01

go 1.14

require (
	github.com/google/go-cmp v0.5.6
	github.com/imdario/mergo v0.3.9 // indirect
	golang.org/x/crypto v0.0.0-20210920023735-84f357641f63
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	k8s.io/api v0.22.5
	k8s.io/apimachinery v0.22.5
	k8s.io/client-go v0.22.5
	knative.dev/hack v0.0.0-20220201013531-82bfca153560
	knative.dev/networking v0.0.0-20220202023132-8ea30455281d
	knative.dev/pkg v0.0.0-20220203020920-51be315ed160
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
