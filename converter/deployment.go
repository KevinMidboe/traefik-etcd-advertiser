package converter

import (
	"errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"os"
)

// NOTE this is really not used, just wanted to convert
// to a kubernetes resource to our packets for practice.
func ServiceToKubernetes(filePath string) (*v1.Service, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	stream, _ := os.ReadFile(filePath)
	// second param (gKV) is top level (GroupVersionKind) w/ group, version & kind
	obj, gKV, _ := decode(stream, nil, nil)

	// handle multiple resources split by ---
	if gKV.Kind == "Service" {
		return obj.(*v1.Service), nil
	}

	return nil, errors.New("Unable to find service resource")
}
