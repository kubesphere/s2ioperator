/*
Copyright 2019 The KubeSphere Authors.

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
	"os"

	"github.com/kubesphere/s2ioperator/pkg/apis"
	"github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	s2iconfig "github.com/kubesphere/s2ioperator/pkg/config"
	"github.com/kubesphere/s2ioperator/pkg/controller"
	"github.com/kubesphere/s2ioperator/pkg/handler"
	"github.com/kubesphere/s2ioperator/pkg/metrics"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func init() {
	// At startup, the default logger in controller runtime is a nil logger.
	// We set it to klogr that is implemented by klog.
	ctrl.SetLogger(klogr.New())
}

func main() {
	var metricsAddr string
	var s2iRunJobTemplatePath string // checkout config/templates/s2irun-template.yaml for example
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&s2iRunJobTemplatePath, "s2irun-job-template", "/etc/template/job.yaml", "the s2irun job template file path")
	flag.Parse()
	log := ctrl.Log.WithName("entrypoint")

	jobTemplateData, err := os.ReadFile(s2iRunJobTemplatePath)
	if err != nil {
		log.Error(err, "failed to read s2irun template file", "filename", s2iRunJobTemplatePath)
		os.Exit(1)
	}
	s2iConfig := &s2iconfig.Config{
		S2IRunJobTemplate: jobTemplateData,
	}

	// Get a config to talk to the apiserver
	log.Info("setting up client for manager")
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "unable to set up client config")
		os.Exit(1)
	}

	// Create a newgo Cmd to provide shared dependencies and start components
	log.Info("setting up manager")
	mgr, err := manager.New(cfg, manager.Options{
		MetricsBindAddress: metricsAddr,
		// We have to set the port to 443 for consistency, because the default port value has been changed
		// from 443 to 9443 after controller-runtime 0.7.0
		// Please see also https://github.com/kubernetes-sigs/controller-runtime/releases/tag/v0.7.0
		Port: 443,
	})
	if err != nil {
		log.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	log.Info("Registering Components.")

	// Setup Scheme for all resources
	log.Info("setting up scheme")
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "unable add APIs to scheme")
		os.Exit(1)
	}

	// Setup all Controllers
	log.Info("Setting up controller")
	if err := controller.AddToManager(mgr, s2iConfig); err != nil {
		log.Error(err, "unable to register controllers to the manager")
		os.Exit(1)
	}

	if err = (&v1alpha1.S2iBuilderTemplate{}).SetupWebhookWithManager(mgr); err != nil {
		log.Error(err, "unable to create webhook", "webhook", "Captain")
		os.Exit(1)
	}

	if err = (&v1alpha1.S2iBuilder{}).SetupWebhookWithManager(mgr); err != nil {
		log.Error(err, "unable to create webhook", "webhook", "Captain")
		os.Exit(1)
	}

	if err = (&v1alpha1.S2iRun{}).SetupWebhookWithManager(mgr); err != nil {
		log.Error(err, "unable to create webhook", "webhook", "Captain")
		os.Exit(1)
	}

	// Set up s2i metrics
	log.Info("start collect s2i metrics")
	go metrics.CollectS2iMetrics(mgr.GetClient())

	// Start webhook handler
	log.Info("start webhook handler")
	go handler.Run(mgr.GetClient())

	//Start the Cmd
	log.Info("Starting the Cmd.")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "unable to run the manager")
		os.Exit(1)
	}
}
