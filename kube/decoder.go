package kube

import (
	"bytes"

	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func ParseKubeObjects(objYaml []byte) ObjectGroup {
	objectGroup := ObjectGroup{}

	metaDecoder := yaml.NewYAMLToJSONDecoder(bytes.NewReader(objYaml))
	objDecoder := yaml.NewYAMLToJSONDecoder(bytes.NewReader(objYaml))

	for {
		var metaObj meta.TypeMeta
		eof := metaDecoder.Decode(&metaObj)
		if eof != nil {
			break
		}

		switch metaObj.Kind {
		case "Deployment":
			deployment := apps.Deployment{}
			objDecoder.Decode(&deployment)
			objectGroup.Deployment = deployment
		case "ConfigMap":
			confMap := core.ConfigMap{}
			objDecoder.Decode(&confMap)
			objectGroup.ConfigMap = confMap
		default:
			//Read current object to advance current reader pointer
			objDecoder.Decode(&meta.TypeMeta{})
		}
	}
	return objectGroup
}
