package bpm_test

import (
	"bytes"
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/JulzDiverse/fusion/bpm"
	"github.com/JulzDiverse/fusion/bpm/bpmfakes"
	"github.com/JulzDiverse/fusion/docker"
)

var _ = Describe("Bpm", func() {

	Context("When a Dockerfile is provided", func() {

		const processName = "opi"

		var (
			bpmer      BPMer
			dockerfile io.Reader
			expected   BPM
			fakeParser *bpmfakes.FakeEntrypointParser
		)

		BeforeEach(func() {

			dockerfile = bytes.NewBufferString(`FROM x

COPY opi /workspace/jobs/opi/bin/

ENTRYPOINT [ "/workspace/jobs/opi/bin/opi", \
	"connect", \
	"--config", \
	"/workspace/jobs/opi/config/opi.yml" \
]`)

			expected = BPM{
				Processes: []Process{
					{
						Name:       processName,
						Executable: "/bin/opi",
						Args: []string{
							"--arg",
							"hello",
						},
						Limits: Limits{
							Memory:    "3G",
							Processes: 10,
							OpenFiles: 100000,
						},
						EphemeralDisk: true,
					},
				},
			}
		})

		JustBeforeEach(func() {
			fakeParser = new(bpmfakes.FakeEntrypointParser)
			bpmer = BPMer{fakeParser}

			fakeParser.ParseDockerfileEntrypointReturns(docker.Entrypoint{
				Executable: "/bin/opi",
				Args: []string{
					"--arg",
					"hello",
				},
			}, nil)
		})

		It("should turn it into a bpm file", func() {
			bpm, err := bpmer.ToBpm(processName, dockerfile)
			Expect(err).ToNot(HaveOccurred())
			Expect(bpm).To(Equal(expected))
		})
	})
})
