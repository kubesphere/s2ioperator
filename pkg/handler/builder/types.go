/*
Copyright 2019 The Kubesphere Authors.

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
	"flag"
	"github.com/golang/glog"
	"net/http"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	Namespace      = "namespace"
	S2iBuilderName = "builder"
)

var KubeClientset client.Client

// Trigger defines the sink resource for processing incoming events.
type Trigger interface {
	Serve(w http.ResponseWriter, r *http.Request)
	ValidateTrigger(eventType string, payload []byte) ([]byte, error)
	Action(eventType string, payload []byte) error
}

type HandlerBuilder struct {
	Pattern string
	Func    http.HandlerFunc
}

func init() {
	var metricsAddr string
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.Parse()

	cfg, err := config.GetConfig()
	if err != nil {
		glog.Error(err, "unable to set up client config")
		os.Exit(1)
	}

	// Create a newgo Cmd to provide shared dependencies and start components
	glog.Info("setting up manager")
	mgr, err := manager.New(cfg, manager.Options{MetricsBindAddress: metricsAddr})

	KubeClientset = mgr.GetClient()
	if err != nil {
		glog.Error(err, "Failed to get the Kubernetes client")
		os.Exit(1)
	}

}

func ClientSets() client.Client {
	return KubeClientset
}
