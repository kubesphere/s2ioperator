package general

import (
	"context"
	guuid "github.com/google/uuid"
	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
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

func (g *Trigger) Action() (err error) {
	namespaceName := types.NamespacedName{
		Namespace: g.Namespace}
	repeat := true

	// really we need check name repeat? This is a very low probability thing.
	for repeat == true {
		s2irunName := s2irunNamePre + guuid.New().String()[:18]
		namespaceName.Name = s2irunName
		repeat, err = g.CheckS2IRunName(namespaceName)
		if err != nil {
			return err
		}
	}

	// create s2irun resource
	s2irun := GenerateNewS2Irun(&namespaceName, g.S2iBuilderName)
	err = g.KubeClientSet.Create(context.TODO(), s2irun)
	if err != nil {
		log.Error(err, "Can not create S2IRun.")
		return err
	}

	return nil
}

func (g *Trigger) CheckS2IRunName(namespaceName types.NamespacedName) (bool, error) {
	// if generate S2IRun name repeat.
	instance := &devopsv1alpha1.S2iRun{}
	err := g.KubeClientSet.Get(context.TODO(), namespaceName, instance)
	if err == nil {
		log.Error(err, "Generate S2IRun name repeat.")
		return true, nil
	}

	// If object not found, continue.
	if errors.IsNotFound(err) {
		return false, nil
	} else {
		return true, err
	}

}

func GenerateNewS2Irun(namespaceName *types.NamespacedName, s2ibuilderName string) *devopsv1alpha1.S2iRun {
	s2irun := &devopsv1alpha1.S2iRun{
		ObjectMeta: v1.ObjectMeta{
			Name:      namespaceName.Name,
			Namespace: namespaceName.Namespace,
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
