package docker

type Entrypoint struct {
	Executable string   `yaml:"executable"`
	Args       []string `yaml:"args"`
}
