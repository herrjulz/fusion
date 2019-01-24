package bpm

import (
	"io"

	"github.com/JulzDiverse/fusion/docker"
)

//go:generate counterfeiter . EntrypointParser
type EntrypointParser interface {
	ParseDockerfileEntrypoint(io.Reader) (docker.Entrypoint, error)
}

type BPMer struct {
	Parser EntrypointParser
}

type BPM struct {
	Processes []Process `yaml:"processes"`
}

type Process struct {
	Name          string   `yaml:"name"`
	Executable    string   `yaml:"executable"`
	Args          []string `yaml:"args"`
	Limits        Limits   `yaml:"limits"`
	EphemeralDisk bool     `yaml:"ephemeral_disk"`
}

type Limits struct {
	Memory    string `yaml:"memory"`
	Processes int    `yaml:"processes"`
	OpenFiles int    `yaml:"open_files"`
}

func (b *BPMer) ToBpm(processName string, dockerfile io.Reader) (BPM, error) {
	entrypoint, err := b.Parser.ParseDockerfileEntrypoint(dockerfile)
	if err != nil {
		return BPM{}, err
	}

	return BPM{
		Processes: []Process{
			{
				Name:       processName,
				Executable: entrypoint.Executable,
				Args:       entrypoint.Args,
				Limits: Limits{
					Memory:    "3G",
					Processes: 10,
					OpenFiles: 100000,
				},
				EphemeralDisk: true,
			},
		},
	}, nil
}
