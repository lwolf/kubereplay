package sharedinformers

// Created by "kubebuilder create resource" for you to start additional shared informers

// SetupKubernetesTypes registers the config for watching Kubernetes types
func (si *SharedInformers) SetupKubernetesTypes() bool {
	return true
}

// StartAdditionalInformers starts watching Deployments
func (si *SharedInformers) StartAdditionalInformers(shutdown <-chan struct{}) {
	go si.KubernetesFactory.Extensions().V1beta1().Deployments().Informer().Run(shutdown)
	go si.KubernetesFactory.Core().V1().ConfigMaps().Informer().Run(shutdown)
}
