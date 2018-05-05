package refinery_test

import (
	"github.com/kubernetes-sigs/kubebuilder/pkg/controller/types"
	kubereplayv1alpha1 "github.com/lwolf/kubereplay/pkg/apis/kubereplay/v1alpha1"
	clientsetv1alpha1 "github.com/lwolf/kubereplay/pkg/client/clientset/versioned/typed/kubereplay/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Refinery controller", func() {
	var instance kubereplayv1alpha1.Refinery
	var expectedKey types.ReconcileKey
	var client clientsetv1alpha1.RefineryInterface

	BeforeEach(func() {
		instance = kubereplayv1alpha1.Refinery{
			Spec: kubereplayv1alpha1.RefinerySpec{
				Storage: &kubereplayv1alpha1.RefineryStorage{
					Stdout: &kubereplayv1alpha1.StdoutSilo{
						Enabled: true,
					},
				},
			},
		}
		instance.Name = "instance-stdout"
		expectedKey = types.ReconcileKey{
			Namespace: "default",
			Name:      "instance-stdout",
		}
	})

	deletePropagation := metav1.DeletePropagationForeground
	AfterEach(func() {
		client.Delete(
			instance.Name,
			&metav1.DeleteOptions{
				PropagationPolicy: &deletePropagation,
			})
	})

	Describe("when creating a new object", func() {
		It("invoke the reconcile method", func() {
			after := make(chan struct{})
			ctrl.AfterReconcile = func(key types.ReconcileKey, err error) {
				defer func() {
					// Recover in case the key is reconciled multiple times
					defer func() { recover() }()
					close(after)
				}()
				defer GinkgoRecover()
				Expect(key).To(Equal(expectedKey))
				Expect(err).ToNot(HaveOccurred())
			}

			// Create the instance
			client = cs.KubereplayV1alpha1().Refineries("default")
			_, err := client.Create(&instance)
			Expect(err).ShouldNot(HaveOccurred())

			// Wait for reconcile to happen
			Eventually(after, "10s", "100ms").Should(BeClosed())

			// INSERT YOUR CODE HERE - test conditions post reconcile
		})
	})
})
