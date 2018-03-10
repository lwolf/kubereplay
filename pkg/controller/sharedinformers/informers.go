package sharedinformers

// EDIT THIS FILE!
// Created by "kubebuilder create resource" for you to start additional shared informers

// SetupKubernetesTypes registers the config for watching Kubernetes types
func (si *SharedInformers) SetupKubernetesTypes() bool {
	return true
}

// StartAdditionalInformers starts watching Deployments
func (si *SharedInformers) StartAdditionalInformers(shutdown <-chan struct{}) {
	// Start specific Kubernetes API informers here.  Note, it is only necessary
	// to start 1 informer for each Kind. (e.g. only 1 Deployment informer)

	// Uncomment this to start listening for Deployment Create / Update / Deletes
	// go si.KubernetesFactory.Apps().V1beta1().Deployments().Informer().Run(shutdown)

	// INSERT YOUR CODE HERE - start additional shared informers
}
