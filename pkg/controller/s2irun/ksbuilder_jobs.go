package s2irun

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	devopsv1alpha1 "github.com/kubesphere/s2ioperator/pkg/apis/devops/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	ConfigDataKey = "data"
)

func (r *ReconcileS2iRun) NewConfigMap(instance *devopsv1alpha1.S2iRun, config *devopsv1alpha1.S2iConfig, template *devopsv1alpha1.UserDefineTemplate) (*corev1.ConfigMap, error) {
	if template != nil {
		t := &devopsv1alpha1.S2iBuilderTemplate{}
		err := r.Get(context.TODO(), types.NamespacedName{Name: template.Name}, t)
		if err != nil {
			return nil, err
		}
		config.BuilderImage = template.BaseImage
	}

	if instance.Spec.NewTag != "" {
		config.Tag = config.ImageName + ":" + instance.Spec.NewTag
	} else {
		config.Tag = config.ImageName + ":" + config.Tag
	}
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	dataMap := make(map[string]string)
	dataMap[ConfigDataKey] = string(data)
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-configmap",
			Namespace: instance.ObjectMeta.Namespace,
		},
		Data: dataMap,
	}
	return configMap, nil
}

func (r *ReconcileS2iRun) GenerateNewJob(instance *devopsv1alpha1.S2iRun) (*batchv1.Job, error) {
	cfgString := "config-data"
	configMapName := instance.Name + "-configmap"
	imageName := os.Getenv("S2IIMAGENAME")
	if imageName == "" {
		return nil, fmt.Errorf("Failed to get s2i-image name, please set the env 'S2IIMAGENAME' ")
	}
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name + "-job",
			Namespace: instance.ObjectMeta.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"job-name": instance.Name + "-job"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "s2irun",
							Image:           imageName,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Env: []corev1.EnvVar{
								{
									Name:  "S2I_CONFIG_PATH",
									Value: "/etc/data/config.json",
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      cfgString,
									ReadOnly:  true,
									MountPath: "/etc/data",
								},
								{
									Name:      "docker-sock",
									MountPath: "/var/run/docker.sock",
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
					Volumes: []corev1.Volume{
						{
							Name: cfgString,
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: configMapName,
									},
									Items: []corev1.KeyToPath{
										{
											Key:  ConfigDataKey,
											Path: "config.json",
										},
									},
								},
							},
						},
						{
							Name: "docker-sock",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{Path: "/var/run/docker.sock"},
							},
						},
					},
				},
			},
			BackoffLimit: &instance.Spec.BackoffLimit,
		},
	}
	if instance.Spec.SecondsAfterFinished > 0 {
		job.Spec.TTLSecondsAfterFinished = &instance.Spec.SecondsAfterFinished
	}
	return job, nil
}
