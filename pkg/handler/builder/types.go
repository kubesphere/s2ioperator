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
	"net/http"
)

const (
	Namespace      = "namespace"
	S2iBuilderName = "builder"
)

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
