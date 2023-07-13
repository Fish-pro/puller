/*
Copyright 2023 The puller Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package builder

import (
	restclient "k8s.io/client-go/rest"
	"k8s.io/controller-manager/pkg/clientbuilder"
	"k8s.io/klog/v2"

	pullerversioned "github.com/puller-io/puller/pkg/generated/clientset/versioned"
)

// ControllerClientBuilder allows you to get clients and configs for application controllers
type ControllerClientBuilder interface {
	clientbuilder.ControllerClientBuilder
	PullerClient(name string) (pullerversioned.Interface, error)
	PullerClientOrDie(name string) pullerversioned.Interface
}

// make sure that PullerControllerClientBuilder implements PullerControllerClientBuilder
var _ PullerControllerClientBuilder = PullerControllerClientBuilder{}

// NewPullerControllerClientBuilder creates a PullerControllerClientBuilder
func NewPullerControllerClientBuilder(config *restclient.Config) ControllerClientBuilder {
	return PullerControllerClientBuilder{
		clientbuilder.SimpleControllerClientBuilder{
			ClientConfig: config,
		},
	}
}

// PullerControllerClientBuilder returns a fixed client with different user agents
type PullerControllerClientBuilder struct {
	clientbuilder.ControllerClientBuilder
}

// PullerClient returns a versioned.Interface built from the ClientBuilder
func (b PullerControllerClientBuilder) PullerClient(name string) (pullerversioned.Interface, error) {
	clientConfig, err := b.Config(name)
	if err != nil {
		return nil, err
	}
	return pullerversioned.NewForConfig(clientConfig)
}

// PullerClientOrDie returns a versioned.interface built from the ClientBuilder with no error.
// If it gets an error getting the client, it will log the error and kill the process it's running in.
func (b PullerControllerClientBuilder) PullerClientOrDie(name string) pullerversioned.Interface {
	client, err := b.PullerClient(name)
	if err != nil {
		klog.Fatal(err)
	}
	return client
}
