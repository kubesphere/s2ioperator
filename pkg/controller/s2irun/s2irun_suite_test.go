package s2irun

import (
	stdlog "log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/kubesphere/s2ioperator/pkg/apis"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	cfg        *rest.Config
	mgr        manager.Manager
	c          client.Client
	recFn      reconcile.Reconciler
	requests   chan reconcile.Request
	stopMgr    chan struct{}
	mgrStopped *sync.WaitGroup
)

func TestS2irun(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "S2irun Suite")
}

var _ = BeforeSuite(func() {
	t := &envtest.Environment{
		CRDDirectoryPaths:        []string{filepath.Join("..", "..", "..", "config", "crds")},
		ControlPlaneStartTimeout: time.Minute * 1,
	}
	apis.AddToScheme(scheme.Scheme)

	var err error
	if cfg, err = t.Start(); err != nil {
		stdlog.Fatal(err)
	}
	mgr, err = manager.New(cfg, manager.Options{})
	Expect(err).NotTo(HaveOccurred())
	c = mgr.GetClient()
	os.Setenv("S2IIMAGENAME", "S2IIMAGENAME/S2IIMAGENAME")

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	recFn, requests = SetupTestReconcile(newReconciler(mgr))
	Expect(add(mgr, recFn)).NotTo(HaveOccurred())
	stopMgr, mgrStopped = StartTestManager(mgr)
})

var _ = AfterSuite(func() {
	os.Unsetenv("S2IIMAGENAME")
	close(stopMgr)
	mgrStopped.Wait()
})

// SetupTestReconcile returns a reconcile.Reconcile implementation that delegates to inner and
// writes the request to requests after Reconcile is finished.
func SetupTestReconcile(inner reconcile.Reconciler) (reconcile.Reconciler, chan reconcile.Request) {
	requests := make(chan reconcile.Request)
	fn := reconcile.Func(func(req reconcile.Request) (reconcile.Result, error) {
		result, err := inner.Reconcile(req)
		requests <- req
		return result, err
	})
	return fn, requests
}

// StartTestManager adds recFn
func StartTestManager(mgr manager.Manager) (chan struct{}, *sync.WaitGroup) {
	stop := make(chan struct{})
	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		Expect(mgr.Start(stop)).NotTo(HaveOccurred())
		wg.Done()
	}()
	return stop, wg
}
