package general

import (
	"context"
	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	log "k8s.io/klog"
	"net/http"
	"path"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func NewTrigger(client client.Client) *Trigger {
	return &Trigger{
		KubeClientSet: client,
	}
}

func (g *Trigger) Serve(w http.ResponseWriter, r *http.Request) {

	reqSecretCode := r.URL.Query().Get("secretCode")

	//example url: host/namespace/buildername
	dir, s2iBuilderName := path.Split(r.URL.Path)
	g.Namespace = path.Base(dir)
	g.S2iBuilderName = s2iBuilderName

	// Authentication
	res, err := g.Authentication(reqSecretCode)
	if err != nil {
		log.Error(err, "Failed to handle event")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !res {
		log.Error(err, "Unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// create resource
	err = g.Action()
	if err != nil {
		log.Error(err, "Failed to handle event")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (g *Trigger) Authentication(reqSecretCode string) (bool, error) {
	s2ibuilder := &devopsv1alpha1.S2iBuilder{}
	namespaceName := types.NamespacedName{
		Name:      g.S2iBuilderName,
		Namespace: g.Namespace}
	err := g.KubeClientSet.Get(context.TODO(), namespaceName, s2ibuilder)
	if err != nil {
		log.Error(err, "Can not get S2IBuilder.")
		return false, err
	}

	if s2ibuilder.Spec.Config.SecretCode == reqSecretCode {
		return true, nil
	} else {
		return false, nil
	}
}

func (g *Trigger) Action() error {

	// generate s2irun resource
	s2irun := g.GenerateNewS2Irun()
	err := g.KubeClientSet.Create(context.TODO(), s2irun)
	if err != nil {
		log.Error(err, "Can not create S2IRun.")
		return err
	}

	return nil
}

func (g *Trigger) GenerateNewS2Irun() *devopsv1alpha1.S2iRun {
	s2irun := &devopsv1alpha1.S2iRun{
		ObjectMeta: v1.ObjectMeta{
			GenerateName: g.S2iBuilderName,
			Namespace:    g.Namespace,
			Annotations: map[string]string{
				"kubesphere.io/creator": defaultCreater,
			},
		},
		Spec: devopsv1alpha1.S2iRunSpec{
			BuilderName: g.S2iBuilderName,
		},
	}

	return s2irun
}
