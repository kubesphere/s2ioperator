package metrics

import (
	"context"
	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/api/core/v1"
	log "k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	"time"
)

var (
	s2iSubsystem = "s2i"

	S2iRunActive = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: s2iSubsystem,
		Name:      "s2irun_active",
		Help:      "Number of s2irun running",
	}, []string{"namespace"})

	S2iRunFailed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: s2iSubsystem,
		Name:      "s2irun_failed",
		Help:      "Number of s2irun failed",
	}, []string{"namespace"})

	S2iRunSucceed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: s2iSubsystem,
		Name:      "s2irun_succeed",
		Help:      "Number of s2irun succeed",
	}, []string{"namespace"})

	S2iRunUnknown = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: s2iSubsystem,
		Name:      "s2irun_unknown",
		Help:      "Number of s2irun which status unknown",
	}, []string{"namespace"})

	S2iBuilderCreated = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: s2iSubsystem,
		Name:      "s2ibuilder_created",
		Help:      "Number of s2ibuilder",
	}, []string{"namespace"})
)

func init() {
	// register the metrics with prometheus registry
	metrics.Registry.MustRegister(S2iRunActive)
	metrics.Registry.MustRegister(S2iRunFailed)
	metrics.Registry.MustRegister(S2iRunSucceed)
	metrics.Registry.MustRegister(S2iRunUnknown)
	metrics.Registry.MustRegister(S2iBuilderCreated)
}

func CollectS2iMetrics(k8sclient client.Client) {
	var err error

	for {
		namespacelist := new(v1.NamespaceList)
		err = k8sclient.List(context.TODO(), namespacelist, client.InNamespace(""))
		if err != nil {
			continue
		}

		for _, namespace := range namespacelist.Items {
			// set s2ibuilder metrics
			err = SetS2iBuilderMetrics(k8sclient, namespace)
			if err != nil {
				log.Error(err)
				continue
			}
			// set s2irun metrics
			err = SetS2iRunMetrics(k8sclient, namespace)
			if err != nil {
				log.Error(err)
				continue
			}

		}
		time.Sleep(15 * time.Second)
	}

}

// Setup s2ibuilder metrics
func SetS2iBuilderMetrics(k8sclient client.Client, namespace v1.Namespace) error {
	s2iBuilderList := new(devopsv1alpha1.S2iBuilderList)
	err := k8sclient.List(context.TODO(), s2iBuilderList, client.InNamespace(namespace.Name))
	if err != nil {
		log.Error(err)
		return err
	}
	S2iBuilderCreated.WithLabelValues(namespace.Name).Set(float64(len(s2iBuilderList.Items)))
	return nil
}

// Setup s2irun metrics
func SetS2iRunMetrics(k8sclient client.Client, namespace v1.Namespace) error {
	s2iRunList := new(devopsv1alpha1.S2iRunList)
	err := k8sclient.List(context.TODO(), s2iRunList, client.InNamespace(namespace.Name))
	if err != nil {
		return err
	}

	var successfulCount = 0
	var runningCount = 0
	var failedCount = 0
	var unknownCount = 0
	for _, s2irun := range s2iRunList.Items {
		switch s2irun.Status.RunState {
		case devopsv1alpha1.Successful:
			successfulCount = successfulCount + 1
		case devopsv1alpha1.Failed:
			failedCount = failedCount + 1
		case devopsv1alpha1.Running:
			runningCount = runningCount + 1
		default:
			unknownCount = unknownCount + 1
		}
	}

	S2iRunSucceed.WithLabelValues(namespace.Name).Set(float64(successfulCount))
	S2iRunActive.WithLabelValues(namespace.Name).Set(float64(runningCount))
	S2iRunFailed.WithLabelValues(namespace.Name).Set(float64(failedCount))
	S2iRunUnknown.WithLabelValues(namespace.Name).Set(float64(unknownCount))
	return nil
}
