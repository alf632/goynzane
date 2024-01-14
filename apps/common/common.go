package common

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func BindMount(src, tgt string) {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil && err.Error() != syscall.ENOTDIR.Error() {
		log.Fatal("cannot create bindmount src folder", src, err)
	}
	err = syscall.Mount(src, tgt, "", syscall.MS_BIND, "")
	if err != nil {
		fmt.Printf("Unable to bind mount %s to %s\t%s\n",
			src, tgt, err)
	}
}

func Run(command ...string) {
	name := command[0]
	args := []string{}
	if len(command) > 1 {
		args = command[1:]
	}
	cmd := exec.Command(name, args...)
	wireIO(cmd)
	if err := cmd.Run(); err != nil {
		fmt.Printf("%v: %v", cmd.Args, err)
	}
}

func RunWithEnv(env, name string, args ...string) {
	cmd := exec.Command(name, args...)
	wireIO(cmd)
	cmd.Env = append(cmd.Env, env)
	if err := cmd.Run(); err != nil {
		fmt.Printf("%v: %v", cmd.Args, err)
	}
}

func Start(name string, args ...string) (*exec.Cmd, context.Context) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, name, args...)
	wireIO(cmd)
	err := cmd.Start()
	if err != nil {
		fmt.Println("cannot start", name)
	}
	return cmd, ctx
}

func StartWithEnv(env, name string, args ...string) (*exec.Cmd, context.Context) {
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, name, args...)
	wireIO(cmd)
	cmd.Env = append(cmd.Env, env)
	err := cmd.Start()
	if err != nil {
		fmt.Println("cannot start", name)
	}
	return cmd, ctx
}

func wireIO(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
}

func MountTmpfs(path string, size int64) error {
	if size < 0 {
		panic("MountTmpfs: size < 0")
	}
	var flags uintptr
	flags = syscall.MS_NOATIME | syscall.MS_SILENT
	flags |= syscall.MS_NODEV | syscall.MS_NOEXEC | syscall.MS_NOSUID
	options := ""
	if size >= 0 {
		options = "size=" + strconv.FormatInt(size, 10)
	}
	err := syscall.Mount("tmpfs", path, "tmpfs", flags, options)
	return os.NewSyscallError("mount", err)
}

func Grep(query, path string) bool {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer file.Close()

	found := false
	scanner := bufio.NewScanner(file)
	for i := 1; scanner.Scan(); i++ {
		if strings.Contains(scanner.Text(), query) {
			found = true
		}
	}
	return found
}
