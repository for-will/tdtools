package main

import (
	"os/exec"
	"testing"
)

func Test_OsCmd(t *testing.T) {
	//os.Stdout
	//
	cmd := exec.Command("C:/cygwin64/bin/mintty.exe",
		//"-e", "/bin/bash",
		"/bin/sh",
		//"-lc", "exec zsh -lc 'date;pwd;ps;who;read'",
		"-lc", "exec bash -lc '/bin/bash -lc \"$(curl -sSL http://172.16.1.220:8081/script/trunk_install.sh)\"'",
		//"/bin/bash -c \"$(curl -sSL http://172.16.1.220:8081/script/trunk_install.sh)\""
		//"-lc", "'pwd;date;sleep 100;'",
	)

	t.Log(cmd.Output())

	//if err := cmd.Run(); err != nil {
	//	t.Fatal(err)
	//}
}
