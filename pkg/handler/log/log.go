package log

import "fmt"

type LogURLGetter interface {
	GetURLOfPodLog(string, string) (string, error)
}

type kubesphereLogGetter struct {
	Version  string
	Protocol string
	URL      string
}

func (k kubesphereLogGetter) GetURLOfPodLog(namespace, pod string) (string, error) {
	return fmt.Sprintf("%s://%s/%s/namespaces/%s/pods/%s?operation=query", k.Protocol, k.URL, k.Version, namespace, pod), nil
}

func GetKubesphereLogger() LogURLGetter {
	return kubesphereLogGetter{
		Version:  "v1alpha2",
		URL:      "ks-apigateway.kubesphere-system.svc/kapis/logging.kubesphere.io",
		Protocol: "http",
	}
}
