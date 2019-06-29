package main

import (
	"encoding/json"
	"fmt"
	"github.com/atotto/clipboard"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// AppConfig -
type AppConfig struct {
	SshArgs []string `json:"sshArgs"`
	Cmd []string `json:"cmd"`
}

var config AppConfig

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	Must(err)

	data, err := ioutil.ReadFile(filepath.Join(dir, "config.json"))
	Must(err)

	err = json.Unmarshal(data, &config)
	Must(err)
}



func main() {
	if len(os.Args) < 5 {
		fmt.Println("Too Few Args: ", os.Args)
		return
	}

	var u *url.URL
	var err error

	if len(os.Args) == 5 {
		fmt.Println("Xshell cmdline")
		u, err = fromXshell()
	} else {
		fmt.Println("SecureCrt cmdline")
		u, err = fromSecureCRT()
	}

	if err != nil {
		fmt.Println(err)
		wait()
	}

	pwd, exist:= u.User.Password()
	if !exist {
		fmt.Println("No Password")
	}

	if err = clipboard.WriteAll(pwd); err != nil {
		fmt.Println(err)
	}

	sshArgs := strings.Join(config.SshArgs, " ")
	ssh := fmt.Sprintf("ssh %s %s@%s -p %s",
		sshArgs, u.User.Username(),u.Hostname(), u.Port())

	fmt.Println("Start Exec: ", u.Fragment)

	args := map[string]string{
		"ssh": ssh,
		"fragment": u.Fragment,
	}

	str := FmtNamedVariable(strings.Join(config.Cmd, " "), args)
	fmt.Println("Full Command: ", str)

	cmds := strings.Split(str, " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		wait()
	}

	fmt.Print("Done!")
}