package spec

import (
	"fmt"
	"strings"

	"github.com/JulzDiverse/fusion/bosh"
	"github.com/JulzDiverse/fusion/kube"
	"github.com/imdario/mergo"
	"github.com/smallfish/simpleyaml"
	"github.com/smallfish/simpleyaml/helper/util"
	"k8s.io/api/core/v1"
)

type Converter struct {
}

func (c *Converter) Convert(objectGroup kube.ObjectGroup) (bosh.Spec, error) {
	deployment := objectGroup.Deployment
	name := objectGroup.Deployment.Name
	configMap := objectGroup.ConfigMap

	templates, properties, err := convertConfigMap(configMap.Data, name)
	if err != nil {
		return bosh.Spec{}, err
	}

	volumes := deployment.Spec.Template.Spec.Volumes
	secretTemplates, secretProps := convertSecrets(volumes, name)

	mergo.Merge(&templates, secretTemplates)
	mergo.Merge(&properties, secretProps)

	return bosh.Spec{
		Name:       name,
		Properties: properties,
		Templates:  templates,
		Packages:   []string{name},
	}, nil
}

func convertConfigMap(data map[string]string, name string) (map[string]string, map[string]bosh.Property, error) {
	properties := map[string]bosh.Property{}
	templates := map[string]string{}

	for filename, content := range data {
		yaml, err := simpleyaml.NewYaml([]byte(content))
		if err != nil {
			return map[string]string{}, map[string]bosh.Property{}, err
		}

		paths, err := util.GetAllPaths(yaml)
		if err != nil {
			return map[string]string{}, map[string]bosh.Property{}, err
		}

		for _, v := range paths {
			v = strings.Replace(v, "/", ".", -1)
			properties[toPath(name, v)] = bosh.Property{Description: "<Not-provided-by-the-user>"}
		}
		templates[fmt.Sprintf("%s.erb", filename)] = fmt.Sprintf("config/%s", filename)
	}
	return templates, properties, nil
}

func convertSecrets(volumes []v1.Volume, name string) (map[string]string, map[string]bosh.Property) {
	properties := map[string]bosh.Property{}
	templates := map[string]string{}

	for _, v := range volumes {
		secret := v.VolumeSource.Secret
		if secret != nil {
			for _, s := range secret.Items {
				properties[toPath(name, s.Key)] = bosh.Property{Description: "<Not-provided-by-the-user>"}
				templates[fmt.Sprintf("%s.erb", s.Path)] = fmt.Sprintf("secret/%s", s.Path)
			}
		}
	}
	return templates, properties
}

func toPath(name, path string) string {
	return fmt.Sprintf("%s.%s", name, path)
}
