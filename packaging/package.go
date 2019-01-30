package packaging

import "fmt"

func CreateScript(binaryName, downloadURL string) string {
	return fmt.Sprintf(`set -e -x

	wget %s -O %[2]s
	cp -a %[2]s ${BOSH_INSTALL_TARGET}
	`, downloadURL, binaryName)
}

func CreateSpec(job, binaryName string) string {
	return fmt.Sprintf(`---
name: %s

files:
  - %s
	`, job, binaryName)
}
