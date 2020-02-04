package github

import (
	"bytes"
	"context"
	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"github.com/kubesphere/s2ioperator/pkg/client/clientset/versioned/scheme"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestValidateTrigger(t *testing.T) {
	as2ib := &devopsv1alpha1.S2iBuilder{
		ObjectMeta: v1.ObjectMeta{
			Name: "s2i-a",
		},
		Spec: devopsv1alpha1.S2iBuilderSpec{
			Config: &devopsv1alpha1.S2iConfig{
				RevisionId:       "branch-a",
				BranchExpression: ".*",
			},
		},
	}
	aPayLoad := []byte(`{"ref": "refs/heads/branch-a"}`)

	bs2ib := &devopsv1alpha1.S2iBuilder{
		ObjectMeta: v1.ObjectMeta{
			Name: "s2i-b",
		},
		Spec: devopsv1alpha1.S2iBuilderSpec{
			Config: &devopsv1alpha1.S2iConfig{
				RevisionId: "branch/b",
			},
		},
	}
	bPayLoad := []byte(`{"ref": "refs/heads/branch/b"}`)

	data := []struct {
		S2ib    *devopsv1alpha1.S2iBuilder
		PayLoad []byte
		Result  bool
	}{
		{S2ib: as2ib, PayLoad: aPayLoad, Result: true},
		{S2ib: bs2ib, PayLoad: bPayLoad, Result: true},
		{S2ib: as2ib, PayLoad: bPayLoad, Result: true},
		{S2ib: bs2ib, PayLoad: aPayLoad, Result: false},
	}

	scheme := scheme.Scheme
	fakeKubeClient := fake.NewFakeClientWithScheme(scheme, as2ib, bs2ib)
	err := fakeKubeClient.Create(context.TODO(), as2ib)
	if err != nil && !errors.IsAlreadyExists(err) {
		t.Fatalf("Ceate resource error %s", err)
	}
	err = fakeKubeClient.Create(context.TODO(), bs2ib)
	if err != nil && !errors.IsAlreadyExists(err) {
		t.Fatalf("Ceate resource error %s", err)
	}
	githubSink := NewGithubSink(fakeKubeClient)

	for _, v := range data {
		githubSink.S2iBuilderName = v.S2ib.Name
		res, err := githubSink.ValidateTrigger(pushEvent, v.PayLoad)
		if v.Result {
			if !bytes.Equal(v.PayLoad, res) {
				t.Fatalf("Get err %s", err)
			}
		} else {
			if err != nil && v.Result == true {
				t.Fatalf("Get err %s", err)
			}
		}
	}
}

func TestAction(t *testing.T) {
	s2ib := &devopsv1alpha1.S2iBuilder{
		ObjectMeta: v1.ObjectMeta{
			Name: "s2i-a",
		},
		Spec: devopsv1alpha1.S2iBuilderSpec{
			Config: &devopsv1alpha1.S2iConfig{
				RevisionId:       "branch-a",
				BranchExpression: ".*",
			},
		},
	}

	aPayLoad := []byte(`{
	"head_commit": {
		"id": "1cb224cd3d4c6490c252b549b0577e9373b18242",
		"tree_id": "f9659cdfc8732b0eceefbcdf0da2665abc29dc95",
		"distinct": true,
		"message": "test1\n\nSigned-off-by: zhuxiaoyang <sunzhu@yunify.com>",
		"timestamp": "2020-02-03T17:40:34+08:00",
		"url": "https://github.com/soulseen/devops-python-sample/commit/1cb224cd3d4c6490c252b549b0577e9373b18242",
		"author": {
			"name": "zhuxiaoyang",
			"email": "sunzhu@yunify.com",
			"username": "soulseen"
		},
		"committer": {
			"name": "zhuxiaoyang",
			"email": "sunzhu@yunify.com",
			"username": "soulseen"
		},
		"added": [],
		"removed": [],
		"modified": [
			"README.md"
		]
	}
}`)

	scheme := scheme.Scheme
	fakeKubeClient := fake.NewFakeClientWithScheme(scheme, s2ib)
	githubSink := NewGithubSink(fakeKubeClient)
	githubSink.S2iBuilderName = s2ib.Name

	err := githubSink.Action(pushEvent, aPayLoad)
	if err != nil {
		t.Fatalf("Get err %s", err)
	}
	res := &devopsv1alpha1.S2iRun{}
	namespacedName := types.NamespacedName{Namespace: "", Name: s2irunNamePre + "1cb22"}
	err = fakeKubeClient.Get(context.TODO(), namespacedName, res)
	if err != nil {
		t.Fatalf("Get err %s", err)
	}
	if res.Name != namespacedName.Name {
		t.Fatalf("The name of s2irun not same with %s ", namespacedName.Name)
	}

	if res.Spec.BuilderName != s2ib.Name {
		t.Fatalf("The BuilderName of s2irun not same with %s ", s2ib.Name)
	}

}
