package bosh

type Spec struct {
	Name       string              `yaml:"name"`
	Templates  map[string]string   `yaml:"templates"`
	Packages   []string            `yaml:"packages"`
	Properties map[string]Property `yaml:"properties"`
}

type Property struct {
	Description string      `yaml:"description"`
	Default     interface{} `yaml:"default"`
}
