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

package handler

import (
	"github.com/emicklei/go-restful"
	"github.com/kubesphere/s2ioperator/pkg/handler/general"
	"github.com/kubesphere/s2ioperator/pkg/handler/github"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"

	log "github.com/golang/glog"
)

func Run(kubeClientset client.Client) {
	container := restful.DefaultContainer

	//register general webhook handler, which can handle any handle request from any server.
	container.Add(general.NewTrigger(kubeClientset).WebService())

	//register  github webhook handler
	container.Add(github.NewTrigger(kubeClientset).WebService())

	log.Info("start listening on localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
