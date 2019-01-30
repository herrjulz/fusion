package monit

import "fmt"

func Create(executable string) string {
	return fmt.Sprintf(`check process %[1]s
  with pidfile /var/vcap/sys/run/bpm/%[1]s/%[1]s.pid
  start program "/var/vcap/jobs/bpm/bin/bpm start %[1]s"
  stop program "/var/vcap/jobs/bpm/bin/bpm stop %[1]s"
  group vcap
	`, executable)
}
