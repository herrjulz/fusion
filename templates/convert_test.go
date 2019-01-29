package templates_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"

	"github.com/JulzDiverse/fusion/kube"
	. "github.com/JulzDiverse/fusion/templates"
)

var _ = Describe("Convert", func() {

	var (
		converter     *Converter
		objectGroup   kube.ObjectGroup
		templateFiles map[string][]byte
		opiConf       string
		err           error
	)

	BeforeEach(func() {
		converter = &Converter{}
		opiConf = `opi:
  prop1: <%= p('babymama.opi.prop1') %>
  prop2: <%= p('babymama.opi.prop2') %>
`
	})

	JustBeforeEach(func() {
		objectGroup = kube.ObjectGroup{
			Deployment: apps.Deployment{
				Spec: apps.DeploymentSpec{
					Template: core.PodTemplateSpec{
						Spec: core.PodSpec{
							Volumes: []core.Volume{
								{
									Name: "whatever",
									VolumeSource: core.VolumeSource{
										Secret: &core.SecretVolumeSource{
											Items: []core.KeyToPath{
												{
													Path: "secret.one",
													Key:  "secret-key-1",
												},
												{
													Path: "secret.two",
													Key:  "secret-key-2",
												},
												{
													Path: "secret.three",
													Key:  "secret-key-3",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			ConfigMap: core.ConfigMap{
				Data: map[string]string{
					"opi.yaml": opiConf,
				},
			},
		}
		objectGroup.Deployment.Name = "babymama"
		templateFiles, err = converter.Convert(objectGroup)
	})

	Context("When a kube object group is provided", func() {
		It("should not error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should generate template files", func() {
			Expect(string(templateFiles["opi.yaml.erb"])).To(Equal(opiConf))
			Expect(string(templateFiles["secret.one.erb"])).To(Equal("<%= p('babymama.secret-key-1') %>"))
			Expect(string(templateFiles["secret.two.erb"])).To(Equal("<%= p('babymama.secret-key-2') %>"))
			Expect(string(templateFiles["secret.three.erb"])).To(Equal("<%= p('babymama.secret-key-3') %>"))
		})
	})

})
