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

package handler

import (
	"github.com/kubesphere/s2ioperator/pkg/handler/github"
	"github.com/kubesphere/s2ioperator/pkg/handler/gitlab"
	"github.com/prometheus/common/log"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/golang/glog"
	"github.com/kubesphere/s2ioperator/pkg/handler/builder"
)

var handlers = []*builder.HandlerBuilder{}

func Run(kubeClientset client.Client) {
	// registry handler type
	handlers = append(handlers, &builder.HandlerBuilder{
		Pattern: "/github/",
		Func:    github.NewGithubSink(kubeClientset).Serve,
	})
	log.Info("registering github webhook")

	handlers = append(handlers, &builder.HandlerBuilder{
		Pattern: "/gitlab/",
		Func:    gitlab.NewGitlabSink(kubeClientset).Serve,
	})
	log.Info("registering gitlab webhook")

	for _, handler := range handlers {
		http.HandleFunc(handler.Pattern, handler.Func)
	}
	glog.Fatal(http.ListenAndServe(":8080", nil))
}
