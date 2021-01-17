module github.com/dennismaxjung/expanded-hostpath-provisioner

go 1.15

require (
	k8s.io/api v0.20.2
	k8s.io/apimachinery v0.20.2
	k8s.io/client-go v0.20.2
	k8s.io/klog/v2 v2.4.0
	sigs.k8s.io/sig-storage-lib-external-provisioner/v6 v6.2.0
)

replace (
	k8s.io/client-go => k8s.io/client-go v0.20.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.20.2
)
