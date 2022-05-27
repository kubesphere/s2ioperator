package docker

import "github.com/docker/distribution/reference"

// ValidateDockerReference does the validation for the image ref
func ValidateDockerReference(ref string) error {
	_, err := reference.Parse(ref)
	return err
}
