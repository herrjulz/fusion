package kube_test

import (
	"bytes"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"

	. "github.com/JulzDiverse/fusion/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Decoder", func() {

	Context("When a kubernetes yaml is provided", func() {

		var (
			err                error
			kubeObjects        io.Reader
			parsedObjects      []interface{}
			expectedDeployment v1.Deployment
			//		expectedStatefulSet v1.StatefulSet
		)

		BeforeEach(func() {
			kubeObjects = bytes.NewBufferString(`---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: eirini
spec:
  replicas: 2
  template:
    metadata:
      labels:
        run: /bin/opi
    spec:
      containers:
      - name: opi
        image: opi
        ports:
        - containerPort: 80
`)
			replicas := int32(2)

			expectedDeployment = v1.Deployment{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Deployment",
					APIVersion: "apps/v1",
				},
				Spec: v1.DeploymentSpec{
					Replicas: &replicas,
					Template: core.PodTemplateSpec{
						Spec: core.PodSpec{
							Containers: []core.Container{
								core.Container{
									Name:  "opi",
									Image: "opi",
									Ports: []core.ContainerPort{
										{ContainerPort: int32(80)},
									},
								},
							},
						},
					},
				},
			}
			expectedDeployment.Name = "eirini"
			expectedDeployment.Spec.Template.Labels = map[string]string{
				"run": "/bin/opi",
			}
		})

		It("should parse the objects in the yaml", func() {
			parsedObjects, err = ParseKubeObjects(kubeObjects)
			Expect(err).ToNot(HaveOccurred())
			Expect(parsedObjects[0].(v1.Deployment)).To(Equal(expectedDeployment))
		})
	})

})
