package entrypoint_test

import (
	"bytes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/JulzDiverse/fusion/docker"
	. "github.com/JulzDiverse/fusion/docker/entrypoint"
)

var _ = Describe("Parser", func() {

	Context("When a Dockerfile is provided", func() {

		var (
			dockerfile *bytes.Buffer
			expected   docker.Entrypoint
			entrypoint docker.Entrypoint
			err        error
		)

		JustBeforeEach(func() {
			entrypoint, err = Parse(dockerfile)
		})

		BeforeEach(func() {
			expected = docker.Entrypoint{
				Executable: "/bin/x",
				Args: []string{
					"--arg",
					"value",
				},
			}
		})

		Context("and it is a simple exec form", func() {
			BeforeEach(func() {
				dockerfile = bytes.NewBufferString(`FROM x
COPY opi /workspace/jobs/opi/bin/
ENTRYPOINT ["/bin/x", "--arg", "value"]`)
			})

			It("should not error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should parse the entrypoint", func() {
				Expect(expected).To(Equal(entrypoint))
			})
		})

		Context("and it is a simple exec form without args", func() {
			BeforeEach(func() {
				dockerfile = bytes.NewBufferString(`FROM x
COPY opi /workspace/jobs/opi/bin/
ENTRYPOINT ["/bin/x"]`)
			})

			It("should not error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should parse the entrypoint", func() {
				Expect(docker.Entrypoint{Executable: "/bin/x", Args: []string{}}).To(Equal(entrypoint))
			})
		})

		Context("and it is a simple exec form without args", func() {
			BeforeEach(func() {
				dockerfile = bytes.NewBufferString(`FROM x
COPY opi /workspace/jobs/opi/bin/
ENTRYPOINT ["/bin/x"]`)
			})

			It("should not error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should parse the entrypoint", func() {
				Expect(docker.Entrypoint{Executable: "/bin/x", Args: []string{}}).To(Equal(entrypoint))
			})
		})

		Context("and it is a simple shell form", func() {
			BeforeEach(func() {
				dockerfile = bytes.NewBufferString(`FROM x
COPY opi /workspace/jobs/opi/bin/
ENTRYPOINT /bin/x --arg value`)
			})

			It("should not error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should parse the entrypoint", func() {
				Expect(expected).To(Equal(entrypoint))
			})
		})

		Context("and it is a shell form without args", func() {
			BeforeEach(func() {
				dockerfile = bytes.NewBufferString(`FROM x
COPY opi /workspace/jobs/opi/bin/
ENTRYPOINT /bin/x`)
			})

			It("should not error", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should parse the entrypoint", func() {
				Expect(docker.Entrypoint{Executable: "/bin/x", Args: []string{}}).To(Equal(entrypoint))
			})
		})

		Context("and it does not contain an entrypoint", func() {
			BeforeEach(func() {
				dockerfile = bytes.NewBufferString(`FROM x
COPY opi /workspace/jobs/opi/bin/`)
			})

			It("should error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring("Could not find entrypoint")))
			})
		})
	})
})
