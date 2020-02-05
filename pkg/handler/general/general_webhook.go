package general

import (
	"context"
	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"github.com/kubesphere/s2ioperator/pkg/handler/builder"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	log "k8s.io/klog"
	"math/rand"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	s2irunNamePre  = "trigger-general-"
	defaultCreater = "auto-trigger"
)

type Trigger struct {
	KubeClientSet  client.Client
	S2iBuilderName string
	Namespace      string
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewGithubSink(client client.Client) *Trigger {
	return &Trigger{
		KubeClientSet: client,
	}
}

func (g *Trigger) Serve(w http.ResponseWriter, r *http.Request) {
	g.Namespace = builder.GetParamInPath(r.URL.Path, builder.Namespace)
	g.S2iBuilderName = builder.GetParamInPath(r.URL.Path, builder.S2iBuilderName)

	err := g.Action()
	if err != nil {
		log.Error(err, "Failed to handle event")
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

func (g *Trigger) Action() error {
	s2irunName := s2irunNamePre + RandStringRunes(5)
	s2irun := GenerateNewS2Irun(s2irunName, g.Namespace, g.S2iBuilderName)
	err := g.KubeClientSet.Create(context.TODO(), s2irun)
	if err != nil {
		log.Error(err, "Can not create S2IRun.")
		return err
	}

	return nil
}

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyz-")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateNewS2Irun(s2irunName, namespace, s2ibuilderName string) *devopsv1alpha1.S2iRun {
	s2irun := &devopsv1alpha1.S2iRun{
		ObjectMeta: v1.ObjectMeta{
			Name:      s2irunName,
			Namespace: namespace,
			Annotations: map[string]string{
				"kubesphere.io/creator": defaultCreater,
			},
		},
		Spec: devopsv1alpha1.S2iRunSpec{
			BuilderName: s2ibuilderName,
		},
	}

	return s2irun
}
