package gitlab

import (
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Trigger struct {
	KubeClientSet  client.Client
	S2iBuilderName string
	Namespace      string
}

func NewGitlabSink(client client.Client) *Trigger {
	return &Trigger{
		KubeClientSet: client,
	}
}

func (g *Trigger) Serve(w http.ResponseWriter, r *http.Request) {

}

func (g *Trigger) ValidateTrigger(eventType string, payload []byte) ([]byte, error) {
	return nil, nil
}

func (g *Trigger) Action(eventType string, payload []byte) error {
	return nil
}
