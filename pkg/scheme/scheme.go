package scheme

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	pullerv1alpha1 "github.com/puller-io/puller/pkg/apis/puller/v1alpha1"
)

// Scheme holds the aggregated Kubernetes's schemes and extended schemes.
var Scheme = runtime.NewScheme()

func init() {
	utilruntime.Must(scheme.AddToScheme(Scheme))
	utilruntime.Must(pullerv1alpha1.AddToScheme(Scheme))
}
