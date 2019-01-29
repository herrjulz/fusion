package templates

import (
	"fmt"
	"log"

	"github.com/JulzDiverse/fusion/kube"
	"github.com/JulzDiverse/goml"
	"github.com/imdario/mergo"
	"k8s.io/api/core/v1"
)

type Converter struct{}

func (c *Converter) Convert(objectGroup kube.ObjectGroup) (map[string][]byte, error) {
	name := objectGroup.Deployment.Name
	data := objectGroup.ConfigMap.Data

	files, err := convertConfigMap(data, name)
	if err != nil {
		log.Fatal(err)
	}

	volumes := objectGroup.Deployment.Spec.Template.Spec.Volumes
	secretTemplates := convertSecrets(volumes, name)
	mergo.Merge(&files, secretTemplates)

	return files, nil
}

func convertConfigMap(data map[string]string, name string) (map[string][]byte, error) {
	files := map[string][]byte{}
	for filename, contents := range data {
		paths, err := goml.GetPaths([]byte(contents))
		if err != nil {
			return map[string][]byte{}, err
		}

		parsed := []byte(contents)
		for _, path := range paths {
			parsed, err = goml.SetInMemory(parsed, path, createRubyProperty(name, path), false)
			if err != nil {
				return map[string][]byte{}, err
			}
		}
		erbFilename := fmt.Sprintf("%s.erb", filename)
		files[erbFilename] = parsed
	}

	return files, nil
}

func convertSecrets(volumes []v1.Volume, name string) map[string][]byte {
	templates := map[string][]byte{}
	for _, v := range volumes {
		secret := v.VolumeSource.Secret
		if secret != nil {
			for _, s := range secret.Items {
				erbFilename := fmt.Sprintf("%s.erb", s.Path)
				templates[erbFilename] = []byte(createRubyProperty(name, s.Key))
			}
		}
	}
	return templates
}

func createRubyProperty(name, path string) string {
	return fmt.Sprintf("<%%= p('%s.%s') %%>", name, path)
}
