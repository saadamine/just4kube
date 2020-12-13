package controllers

import (
	"context"

	"github.com/go-logr/zapr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	urlv1alpha1 "gytigyg.io/api/v1alpha1"
)

var _ = Describe("Friend controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		FriendName      = "friend-sample"
		FriendNamespace = "default"
		Uri             = "https://url.gytigyg.io/amine"
	)
	// configMapName := []string{"configmap"}
	Context("When updating Friend Status", func() {
		It("Should update Friend Status.Active to success when new Friend are created", func() {
			ctx := context.Background()
			friend := &urlv1alpha1.Friend{
				// TypeMeta: metav1.TypeMeta{
				//     APIVersion: "batch.tutorial.kubebuilder.io/v1",
				//     Kind:       "Friend",
				// },
				ObjectMeta: metav1.ObjectMeta{
					Name:      FriendName,
					Namespace: FriendNamespace,
				},
				Spec: urlv1alpha1.FriendSpec{
					Uri: Uri,
					// ConfigMapName: configMapName,
				},
			}
			Expect(k8sClient.Create(ctx, friend)).Should(Succeed())

			name := types.NamespacedName{Namespace: friend.Namespace, Name: friend.Name}

			zapLog, _ := zap.NewDevelopment()
			reconcile := &FriendReconciler{
				Client: k8sClient,
				Log:    zapr.NewLogger(zapLog),
				Scheme: scheme.Scheme,
			}
			Expect(reconcile.Reconcile(ctrl.Request{NamespacedName: name})).To(Equal(ctrl.Result{}))

			response := urlv1alpha1.Friend{}
			Expect(k8sClient.Get(ctx, name, &response)).To(Succeed())
			Expect("Success").To(Equal(response.Status.Active), "Expected reconcile to change the status to Success")
		})
	})
})
