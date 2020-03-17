package general

import (
	"context"
	"github.com/emicklei/go-restful"
	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	log "k8s.io/klog"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
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

func (g *Trigger) Serve(request *restful.Request, response *restful.Response) {

	reqSecretCode := request.QueryParameter("secretCode")
	g.S2iBuilderName = request.PathParameter("s2ibuilder")
	g.Namespace = request.PathParameter("namespace")

	// Authentication
	res, err := g.Authentication(reqSecretCode)
	if err != nil {
		log.Error(err, "Failed to handle event")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !res {
		log.Error(err, "Unauthorized")
		response.WriteHeader(http.StatusUnauthorized)
		return
	}

	// create resource
	err = g.Action()
	if err != nil {
		log.Error(err, "Failed to handle event")
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusCreated)
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

// do something when handler be triggered.
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

// generate S2Irun yaml.
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
