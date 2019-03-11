package e2e_test

import (
	"fmt"
	"log"
	"os/exec"
	"path"
	"runtime"
	"testing"
	"time"

	"github.com/kubesphere/s2ioperator/pkg/apis"
	"github.com/kubesphere/s2ioperator/pkg/util/e2eutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	testClient client.Client
	cfg        *rest.Config
	workspace  string
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}

var _ = BeforeSuite(func() {
	//install deploy
	workspace = getWorkspace() + "/.."
	cfg, err := config.GetConfig()
	Expect(err).ShouldNot(HaveOccurred(), "Error reading kubeconfig")
	apis.AddToScheme(scheme.Scheme)
	c, err := client.New(cfg, client.Options{})
	Expect(err).NotTo(HaveOccurred(), "Error in creating client")
	testClient = c
	//install deployment
	cmd := exec.Command("kubectl", "apply", "-f", workspace+"/deploy/s2ioperator.yaml")
	bytes, err := cmd.CombinedOutput()
	Expect(err).ShouldNot(HaveOccurred())
	log.Println(string(bytes))
	//waiting for controller up
	err = e2eutil.WaitForController(c, "devops-test", "controller-manager", 15*time.Second, 2*time.Minute)
	Expect(err).ShouldNot(HaveOccurred(), "timeout waiting for controller up: %s\n", err)
	fmt.Fprintf(GinkgoWriter, "Controller is up now")
})

var _ = AfterSuite(func() {
	cmd := exec.Command("kubectl", "delete", "-f", workspace+"/deploy/s2ioperator.yaml")
	Expect(cmd.Run()).ShouldNot(HaveOccurred())
})

func getWorkspace() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename)
}
