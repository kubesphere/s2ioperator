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

package defaultserver

import (
	"fmt"

	s2iBuilderMutating "github.com/kubesphere/s2ioperator/pkg/webhook/default_server/s2ibuilder/mutating"
	s2iBuilderValidating "github.com/kubesphere/s2ioperator/pkg/webhook/default_server/s2ibuilder/validating"
	s2iBuilderTemplateValidating "github.com/kubesphere/s2ioperator/pkg/webhook/default_server/s2ibuildertemplate/validating"
	s2iRunValidating "github.com/kubesphere/s2ioperator/pkg/webhook/default_server/s2irun/validating"
)

func init() {
	for k, v := range s2iBuilderValidating.Builders {
		_, found := builderMap[k]
		if found {
			log.V(1).Info(fmt.Sprintf(
				"conflicting webhook builder names in builder map: %v", k))
		}
		builderMap[k] = v
	}
	for k, v := range s2iBuilderValidating.HandlerMap {
		_, found := HandlerMap[k]
		if found {
			log.V(1).Info(fmt.Sprintf(
				"conflicting webhook builder names in handler map: %v", k))
		}
		_, found = builderMap[k]
		if !found {
			log.V(1).Info(fmt.Sprintf(
				"can't find webhook builder name %q in builder map", k))
			continue
		}
		HandlerMap[k] = v
	}
	for k, v := range s2iBuilderTemplateValidating.Builders {
		_, found := builderMap[k]
		if found {
			log.V(1).Info(fmt.Sprintf(
				"conflicting webhook builder names in builder map: %v", k))
		}
		builderMap[k] = v
	}

	for k, v := range s2iBuilderTemplateValidating.HandlerMap {
		_, found := HandlerMap[k]
		if found {
			log.V(1).Info(fmt.Sprintf(
				"conflicting webhook builder names in handler map: %v", k))
		}
		_, found = builderMap[k]
		if !found {
			log.V(1).Info(fmt.Sprintf(
				"can't find webhook builder name %q in builder map", k))
			continue
		}
		HandlerMap[k] = v
	}

	for k, v := range s2iRunValidating.Builders {
		_, found := builderMap[k]
		if found {
			log.V(1).Info(fmt.Sprintf(
				"conflicting webhook builder names in builder map: %v", k))
		}
		builderMap[k] = v
	}

	for k, v := range s2iRunValidating.HandlerMap {
		_, found := HandlerMap[k]
		if found {
			log.V(1).Info(fmt.Sprintf(
				"conflicting webhook builder names in handler map: %v", k))
		}
		_, found = builderMap[k]
		if !found {
			log.V(1).Info(fmt.Sprintf(
				"can't find webhook builder name %q in builder map", k))
			continue
		}
		HandlerMap[k] = v
	}

	for k, v := range s2iBuilderMutating.Builders {
		_, found := builderMap[k]
		if found {
			log.V(1).Info(fmt.Sprintf(
				"conflicting webhook builder names in builder map: %v", k))
		}
		builderMap[k] = v
	}

	for k, v := range s2iBuilderMutating.HandlerMap {
		_, found := HandlerMap[k]
		if found {
			log.V(1).Info(fmt.Sprintf(
				"conflicting webhook builder names in handler map: %v", k))
		}
		_, found = builderMap[k]
		if !found {
			log.V(1).Info(fmt.Sprintf(
				"can't find webhook builder name %q in builder map", k))
			continue
		}
		HandlerMap[k] = v
	}
}
