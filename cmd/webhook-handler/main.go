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
package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/kubesphere/s2ioperator/pkg/handler"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func main() {
	flag.Parse()
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

	if err != nil {
		glog.Error(err, "Failed to get the Kubernetes client")
		os.Exit(1)
	}

	handler.Run(mgr.GetClient())
}
