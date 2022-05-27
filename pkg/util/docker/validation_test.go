package docker

import "testing"

func Test_validateDockerReference(t *testing.T) {
	type args struct {
		ref string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{{
		name: "without image tag",
		args: args{
			ref: "kubesphere/ks-devops",
		},
		wantErr: false,
	}, {
		name: "with a valid image tag",
		args: args{
			ref: "kubesphere/ks-devops:latest",
		},
		wantErr: false,
	}, {
		name: "with image registry server address",
		args: args{
			ref: "docker.io/kubesphere/ks-devops:latest",
		},
		wantErr: false,
	}, {
		name: "a fake image name",
		args: args{
			ref: "fake-name",
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateDockerReference(tt.args.ref); (err != nil) != tt.wantErr {
				t.Errorf("validateDockerReference() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
