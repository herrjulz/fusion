package kube_test

import (
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
			kubeObjects        string
			expectedDeployment v1.Deployment
			expectedConfigMap  core.ConfigMap
		)

		BeforeEach(func() {
			kubeObjects = `---
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
---
apiVersion: v1
kind: Service
metadata:
  name: "eirini-opi"
spec:
  externalIPs: 1.2.3.4
  ports:
    - port: 8085
      protocol: TCP
      name: opi
  selector:
    name: "eirini"
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: "eirini"
data:
  opi.yaml: |
    opi:
      kube_namespace: mars
      nats_password: deeznatz
`
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

			expectedConfigMap = core.ConfigMap{
				Data: map[string]string{
					"opi.yaml": `opi:
  kube_namespace: mars
  nats_password: deeznatz
`,
				},
			}
			expectedConfigMap.Name = "eirini"
			expectedConfigMap.Kind = "ConfigMap"
			expectedConfigMap.APIVersion = "v1"
		})

		It("should parse the deployment", func() {
			objectGroup := ParseKubeObjects([]byte(kubeObjects))
			Expect(objectGroup.Deployment).To(Equal(expectedDeployment))
		})

		FIt("should parse the configMap", func() {
			objectGroup := ParseKubeObjects([]byte(kubeObjects))
			Expect(objectGroup.ConfigMap).To(Equal(expectedConfigMap))
		})
	})

})
