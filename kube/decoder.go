package kube

import (
	"io"

	"k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func ParseKubeObjects(objYaml io.Reader) ([]interface{}, error) {
	objects := []interface{}{}
	decoder := yaml.NewYAMLToJSONDecoder(objYaml)

	deployment := v1.Deployment{}
	decoder.Decode(&deployment)
	objects = append(objects, deployment)
	return objects, nil
}
