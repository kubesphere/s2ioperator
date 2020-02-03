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

package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"github.com/kubesphere/s2ioperator/pkg/handler/builder"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	log "k8s.io/klog"
	"net/http"
	"reflect"
	"regexp"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

const (
	s2irunNamePre    = "trigger-github-"
	s2irunCreatorPre = "trigger-"
	pushEvent        = "PushEvent"
)

type Trigger struct {
	KubeClientSet  client.Client
	S2iBuilderName string
	Namespace      string
}

func NewGithubSink(client client.Client) *Trigger {
	return &Trigger{
		KubeClientSet: client,
	}
}

func (g *Trigger) Serve(w http.ResponseWriter, r *http.Request) {
	g.Namespace = r.URL.Query().Get(builder.Namespace)
	g.S2iBuilderName = r.URL.Query().Get(builder.S2iBuilderName)

	eventType := github.WebHookType(r)
	// Currently only accepting json payloads.
	eventPayload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error reading event body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// validate payload
	payload, err := g.ValidateTrigger(eventType, eventPayload)
	if err != nil {
		log.Error("Failed to validate event")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = g.Action(eventType, payload)
	if err != nil {
		log.Error(err, "Failed to handle event")
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
	log.Infof("Github handing event with S2IBuilder name %s in namespace %s", g.S2iBuilderName, g.Namespace)

}

func (g *Trigger) ValidateTrigger(eventType string, payload []byte) ([]byte, error) {
	instance := &devopsv1alpha1.S2iBuilder{}
	namespacedName := &types.NamespacedName{Namespace: g.Namespace, Name: g.S2iBuilderName}
	err := g.KubeClientSet.Get(context.TODO(), *namespacedName, instance)
	if err != nil {
		log.Errorf("Failed to get S2IBuilder: %s, in namespace %s", g.S2iBuilderName, g.Namespace)
		return nil, err
	}

	// Check if the event type is in the allow-list, Now just support push event.
	if eventType != pushEvent {
		return nil, fmt.Errorf("not support event type %s", eventType)
	}
	if instance.Spec.Config.AllowedEvents != nil {
		isAllowed := false
		for _, allowedEvent := range instance.Spec.Config.AllowedEvents {
			if eventType == allowedEvent {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			return nil, fmt.Errorf("not support event type %s", eventType)
		}
	}
	// Can not get branch name directly.
	event, err := github.ParseWebHook(eventType, payload)
	pushEvent := event.(github.PushEvent)
	gitref := pushEvent.Ref
	branchName := strings.Split(*gitref,"/tags/")[1]
	if instance.Spec.Config.BranchExpression != "" {
		match, err := regexp.MatchString(instance.Spec.Config.BranchExpression, branchName)
		if err != nil {
			log.Error("Failed to MatchString with BranchName %s by Expression %s", instance.Spec.Config.BranchExpression, branchName)
			return nil, err
		}

		if !match {
			return nil, fmt.Errorf("branch %s is not matched", branchName)
		}
	}

	return payload, nil
}

func (g *Trigger) Action(eventType string, payload []byte) (err error) {
	event, err := github.ParseWebHook(eventType, payload)
	t := reflect.TypeOf(event)
	switch t.Name() {
	case pushEvent:
		err = g.actionWithPushEvent(event.(github.PushEvent))
	case "PullRequestEvent":
		err = g.actionWithPullRequestEvent(event.(github.PullRequestEvent))
	default:
		log.Infof("Can not do any action with event type %s", t.Name())
	}

	return err
}

func (g *Trigger) actionWithPushEvent(event github.PushEvent) error {
	revisionId := event.HeadCommit.ID
	creater := s2irunCreatorPre + *event.HeadCommit.Committer.Name
	s2irunName := s2irunNamePre + string([]byte(*revisionId)[0:5])

	s2irun := GenerateNewS2Irun(s2irunName, g.Namespace, creater, g.S2iBuilderName, *revisionId)
	err := g.KubeClientSet.Create(context.TODO(), s2irun)
	if err != nil {
		log.Error(err, "Can not create S2IRun.")
		return err
	}
	return nil
}

func GenerateNewS2Irun(s2irunName, namespace, creator, s2ibuilderName, revisionId string) *devopsv1alpha1.S2iRun {
	s2irun := &devopsv1alpha1.S2iRun{
		ObjectMeta: v1.ObjectMeta{
			Name:      s2irunName,
			Namespace: namespace,
			Annotations: map[string]string{
				"kubesphere.io/creator": creator,
			},
		},
		Spec: devopsv1alpha1.S2iRunSpec{
			BuilderName:   s2ibuilderName,
			NewRevisionId: revisionId,
		},
	}

	return s2irun
}

func (g *Trigger) actionWithPullRequestEvent(github.PullRequestEvent) error {
	log.Info("Can not do any action with event type PullRequest")
	return nil
}
