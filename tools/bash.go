package tools

import "os/exec"

func Bash(cmd string) (out string, exitcode int) {
	cmdobj := exec.Command("bash", "-c", cmd)
	output, err := cmdobj.CombinedOutput()
	if err != nil {
		// Get the exitcode of the output
		if ins, ok := err.(*exec.ExitError); ok {
			out = string(output)
			exitcode = ins.ExitCode()
			return out, exitcode
		}
	}
	return string(output), 0
}
