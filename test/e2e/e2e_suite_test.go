package e2e_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/types"

	"github.com/kubesphere/s2ioperator/pkg/apis"
	"github.com/kubesphere/s2ioperator/pkg/util/e2eutil"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	testClient    client.Client
	cfg           *rest.Config
	workspace     string
	testNamespace string
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2e Suite")
}

var _ = BeforeSuite(func() {
	//install deploy
	//install operator is writing in the `make debug`, maybe we should write here to decouple
	testNamespace = os.Getenv("TEST_NS")
	Expect(testNamespace).ShouldNot(BeEmpty())
	workspace = getWorkspace() + "/../.."
	cfg, err := config.GetConfig()
	Expect(err).ShouldNot(HaveOccurred(), "Error reading kubeconfig")
	apis.AddToScheme(scheme.Scheme)
	c, err := client.New(cfg, client.Options{})
	Expect(err).NotTo(HaveOccurred(), "Error in creating client")
	testClient = c
	//waiting for controller up
	err = e2eutil.WaitForController(c, testNamespace, "s2ioperator", 5*time.Second, 2*time.Minute)
	Expect(err).ShouldNot(HaveOccurred(), "timeout waiting for controller up: %s\n", err)
	//waiting for webhook
	Eventually(func() error {
		service := &corev1.Service{}
		return c.Get(context.TODO(), types.NamespacedName{Namespace: testNamespace, Name: "webhook-server-service"}, service)
	}, time.Minute*1, time.Second*2).Should(Succeed())

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
