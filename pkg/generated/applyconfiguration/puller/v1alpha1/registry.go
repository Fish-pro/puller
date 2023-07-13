/*
Copyright The Kubernetes Authors.

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
// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// RegistryApplyConfiguration represents an declarative configuration of the Registry type for use
// with apply.
type RegistryApplyConfiguration struct {
	Server   *string `json:"server,omitempty"`
	Username *string `json:"username,omitempty"`
	Password *string `json:"password,omitempty"`
	Email    *string `json:"email,omitempty"`
	Auth     *string `json:"auth,omitempty"`
}

// RegistryApplyConfiguration constructs an declarative configuration of the Registry type for use with
// apply.
func Registry() *RegistryApplyConfiguration {
	return &RegistryApplyConfiguration{}
}

// WithServer sets the Server field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Server field is set to the value of the last call.
func (b *RegistryApplyConfiguration) WithServer(value string) *RegistryApplyConfiguration {
	b.Server = &value
	return b
}

// WithUsername sets the Username field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Username field is set to the value of the last call.
func (b *RegistryApplyConfiguration) WithUsername(value string) *RegistryApplyConfiguration {
	b.Username = &value
	return b
}

// WithPassword sets the Password field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Password field is set to the value of the last call.
func (b *RegistryApplyConfiguration) WithPassword(value string) *RegistryApplyConfiguration {
	b.Password = &value
	return b
}

// WithEmail sets the Email field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Email field is set to the value of the last call.
func (b *RegistryApplyConfiguration) WithEmail(value string) *RegistryApplyConfiguration {
	b.Email = &value
	return b
}

// WithAuth sets the Auth field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Auth field is set to the value of the last call.
func (b *RegistryApplyConfiguration) WithAuth(value string) *RegistryApplyConfiguration {
	b.Auth = &value
	return b
}