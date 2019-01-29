package kube

import (
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
)

type ObjectGroup struct {
	Deployment apps.Deployment
	ConfigMap  core.ConfigMap
}
