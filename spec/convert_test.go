package spec_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"

	"github.com/JulzDiverse/fusion/bosh"
	"github.com/JulzDiverse/fusion/kube"
	. "github.com/JulzDiverse/fusion/spec"
)

var _ = Describe("Convert", func() {

	Context("When kube objects are provided", func() {

		var (
			converter   *Converter
			objectGroup kube.ObjectGroup
			spec        bosh.Spec
			err         error
			opiConfig   string
		)

		BeforeEach(func() {
			opiConfig = `
opi:
  kube_config: my-conf.yml
  kube_endpoint: kube.end-my-point.com
  nats_ip: 1.2.3.4
  api_endpoint: api.cc.whatever
`
		})

		JustBeforeEach(func() {
			converter = &Converter{}
			objectGroup = kube.ObjectGroup{
				Deployment: v1.Deployment{
					Spec: v1.DeploymentSpec{
						Template: core.PodTemplateSpec{
							Spec: core.PodSpec{
								Volumes: []core.Volume{
									{
										Name: "whatever",
										VolumeSource: core.VolumeSource{
											Secret: &core.SecretVolumeSource{
												Items: []core.KeyToPath{
													{
														Key:  "secret-key-1",
														Path: "secret.one",
													},
													{
														Key:  "secret-key-2",
														Path: "secret.two",
													},
													{
														Key:  "secret-key-3",
														Path: "secret.three",
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
						"opi.yaml": opiConfig,
					},
				},
			}

			objectGroup.Deployment.Name = "babymama"
			spec, err = converter.Convert(objectGroup)
		})

		It("should not fail", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should be converted into a bosh-release job spec object", func() {
			Expect(spec).To(Equal(bosh.Spec{
				Name: "babymama",
				Properties: map[string]bosh.Property{
					"babymama.opi.kube_config":   bosh.Property{Description: "<Not-provided-by-the-user>"},
					"babymama.opi.kube_endpoint": bosh.Property{Description: "<Not-provided-by-the-user>"},
					"babymama.opi.nats_ip":       bosh.Property{Description: "<Not-provided-by-the-user>"},
					"babymama.opi.api_endpoint":  bosh.Property{Description: "<Not-provided-by-the-user>"},
					"babymama.secret-key-1":      bosh.Property{Description: "<Not-provided-by-the-user>"},
					"babymama.secret-key-2":      bosh.Property{Description: "<Not-provided-by-the-user>"},
					"babymama.secret-key-3":      bosh.Property{Description: "<Not-provided-by-the-user>"},
				},
				Packages: []string{
					"babymama",
				},
				Templates: map[string]string{
					"opi.yaml.erb":     "config/opi.yaml",
					"secret.one.erb":   "secret/secret.one",
					"secret.two.erb":   "secret/secret.two",
					"secret.three.erb": "secret/secret.three",
				},
			}))
		})

		Context("and the ConfigMap does not contain valid YAML", func() {
			BeforeEach(func() {
				opiConfig = `
			this-is-not-yaml
`
			})

			It("should fail", func() {
				Expect(err).To(HaveOccurred())
			})

		})
	})
})
