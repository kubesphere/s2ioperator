package validating

import (
	"strings"

	"github.com/docker/distribution/reference"
	api "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"github.com/kubesphere/s2ioperator/pkg/errors"
)

// ValidateConfig returns a list of error from validation.
func ValidateConfig(config *api.S2iConfig, fromTemplate bool) []error {
	allErrs := make([]error, 0)
	if !config.IsBinaryURL && len(config.SourceURL) == 0 {
		allErrs = append(allErrs, errors.NewFieldRequired("sourceUrl"))
	}
	if !fromTemplate && len(config.BuilderImage) == 0 {
		allErrs = append(allErrs, errors.NewFieldRequired("builderImage"))
	}
	switch config.BuilderPullPolicy {
	case api.PullNever, api.PullAlways, api.PullIfNotPresent:
	default:
		allErrs = append(allErrs, errors.NewFieldInvalidValue("builderPullPolicy"))
	}
	if config.DockerNetworkMode != "" && !validateDockerNetworkMode(config.DockerNetworkMode) {
		allErrs = append(allErrs, errors.NewFieldInvalidValue("dockerNetworkMode"))
	}
	if config.Labels != nil {
		for k := range config.Labels {
			if len(k) == 0 {
				allErrs = append(allErrs, errors.NewFieldInvalidValue("labels"))
			}
		}
	}
	if config.BuilderImage != "" {
		if err := validateDockerReference(config.BuilderImage); err != nil {
			allErrs = append(allErrs, errors.NewFieldInvalidValueWithReason("builderImage", err.Error()))
		}
	}
	if config.RuntimeAuthentication != nil {
		if config.RuntimeAuthentication.SecretRef == nil {
			if config.RuntimeAuthentication.Username == "" && config.RuntimeAuthentication.Password == "" {
				allErrs = append(allErrs, errors.NewFieldRequired("RuntimeAuthentication username|password / secretRef"))
			}
		}
	}
	if config.IncrementalAuthentication != nil {
		if config.IncrementalAuthentication.SecretRef == nil {
			if config.IncrementalAuthentication.Username == "" && config.IncrementalAuthentication.Password == "" {
				allErrs = append(allErrs, errors.NewFieldRequired("IncrementalAuthentication username|password / secretRef"))
			}
		}
	}
	if config.PullAuthentication != nil {
		if config.PullAuthentication.SecretRef == nil {
			if config.PullAuthentication.Username == "" && config.PullAuthentication.Password == "" {
				allErrs = append(allErrs, errors.NewFieldRequired("PullAuthentication username|password / secretRef"))
			}
		}
	}
	if config.PushAuthentication != nil {
		if config.PushAuthentication.SecretRef == nil {
			if config.PushAuthentication.Username == "" && config.PushAuthentication.Password == "" {
				allErrs = append(allErrs, errors.NewFieldRequired("PushAuthentication username|password / secretRef"))
			}
		}
	}
	return allErrs
}

func validateDockerReference(ref string) error {
	_, err := reference.Parse(ref)
	return err
}

// validateDockerNetworkMode checks wether the network mode conforms to the docker remote API specification (v1.19)
// Supported values are: bridge, host, container:<name|id>, and netns:/proc/<pid>/ns/net
func validateDockerNetworkMode(mode api.DockerNetworkMode) bool {
	switch mode {
	case api.DockerNetworkModeBridge, api.DockerNetworkModeHost:
		return true
	}
	if strings.HasPrefix(string(mode), api.DockerNetworkModeContainerPrefix) {
		return true
	}
	if strings.HasPrefix(string(mode), api.DockerNetworkModeNetworkNamespacePrefix) {
		return true
	}
	return false
}

func ValidateParameter(user, tmpt []api.Parameter) []error {
	findParameter := func(name string, ps []api.Parameter) int {
		for index, v := range ps {
			if v.Key == name {
				return index
			}
		}
		return -1
	}
	allErrs := make([]error, 0)
	for _, v := range tmpt {
		index := findParameter(v.Key, user)
		if v.Required && (index == -1 || user[index].Value == "") {
			allErrs = append(allErrs, errors.NewFieldRequired("Parameter:"+v.Key))
		}
	}
	return allErrs
}
