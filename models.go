package fusion

type Config struct {
	BoshReleaseDir     string `yaml:"bosh_release_dir"`
	KubeDeploymentSpec string `yaml:"kube_deployment_spec"`
	Dockerfile         string `yaml:"dockerfile"`
	BinaryDownloadURL  string `yaml:"binary_download_url"`
}
