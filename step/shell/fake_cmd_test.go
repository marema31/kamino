package shell // We will extend the type only in test environment

import (
	"fmt"
	"os"
) // FakeCmd modify the step.cmd member to not start the original command but a Go function
func (st *Step) FakeCmd(returnCode int, addedEnv []string) {
	st.cmd.Args = append([]string{"-test.run=TestHelperProcess", "--", st.cmd.Path}, st.cmd.Args...)
	st.cmd.Path = os.Args[0] // The go test command
	st.cmd.Env = append(st.cmd.Env, st.cmd.Env...)
	st.cmd.Env = append(st.cmd.Env, "GO_WANT_HELPER_PROCESS=1")
	if returnCode != 0 {
		st.cmd.Env = append(st.cmd.Env, fmt.Sprintf("TESTHELPEREXIT=%d", returnCode))
	}
}
