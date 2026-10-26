package connect

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
)

const (
	tunnelReadyTimeout  = 15 * time.Second
	tunnelPollInterval  = 100 * time.Millisecond
	tunnelDialTimeout   = 100 * time.Millisecond
	tunnelLocalBindAddr = "127.0.0.1"
	tunnelBastionUser   = "bastion"
)

// Tunnel is a child process that forwards a local TCP port to a remote
// address. It is started before a connect session and stopped when it ends.
type Tunnel struct {
	cmd       *exec.Cmd
	localPort uint16
}

// LocalAddr returns the loopback address (host:port) that reaches the
// tunneled remote address.
func (t *Tunnel) LocalAddr() string {
	return net.JoinHostPort(tunnelLocalBindAddr, strconv.Itoa(int(t.localPort)))
}

// LocalPort returns the local TCP port that reaches the tunneled address.
func (t *Tunnel) LocalPort() uint16 {
	return t.localPort
}

// sshLocalForwardArgs builds the ssh command line that forwards a local port
// to host:port through the SSH bastion of a Scaleway VPC gateway:
//
//	ssh -N -o ExitOnForwardFailure=yes -L 127.0.0.1:<local>:<host>:<port> -p <bastion-port> bastion@<gateway-ip>
func sshLocalForwardArgs(
	localPort uint16,
	host string,
	port, bastionPort uint32,
	gatewayIP string,
) []string {
	return []string{
		"-N",
		"-o", "ExitOnForwardFailure=yes",
		"-L", fmt.Sprintf("%s:%d:%s:%d", tunnelLocalBindAddr, localPort, host, port),
		"-p", strconv.Itoa(int(bastionPort)),
		tunnelBastionUser + "@" + gatewayIP,
	}
}

// SSHLocalForward starts an SSH tunnel through the SSH bastion of the given
// VPC gateway that forwards a fresh local port to host:port. The tunnel is
// ready when the local port accepts connections; Stop must be called when
// the session ends. ssh output (auth failures, host key prompts, ...) is
// written to stderr.
func SSHLocalForward(
	stderr io.Writer,
	gatewayIP string,
	bastionPort uint32,
	host string,
	port uint32,
) (*Tunnel, error) {
	localPort, err := freeLocalPort()
	if err != nil {
		return nil, &core.CliError{
			Err:  fmt.Errorf("failed to allocate a local port for the tunnel: %w", err),
			Hint: "the system seems to be out of local TCP ports: close idle connections or applications and retry",
		}
	}

	t := &Tunnel{
		localPort: localPort,
		cmd: exec.Command(
			"ssh",
			sshLocalForwardArgs(localPort, host, port, bastionPort, gatewayIP)...),
	}
	t.cmd.Stderr = stderr
	// Share the terminal with the connect command: ssh may need to prompt
	// for the host key or a key passphrase. With -N ssh does not read stdin
	// once authenticated, so this does not interfere with the connect tool.
	t.cmd.Stdin = os.Stdin

	if err := t.cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start the ssh tunnel: %w", err)
	}

	// The waiter goroutine owns cmd.Wait; Stop only kills the process. The
	// channel is buffered so the goroutine never leaks.
	exitCh := make(chan error, 1)
	go func() {
		exitCh <- t.cmd.Wait()
	}()

	if err := t.WaitForReady(tunnelReadyTimeout, exitCh); err != nil {
		_ = t.cmd.Process.Kill()

		return nil, err
	}

	return t, nil
}

// WaitForReady polls the local port until it accepts connections. exitCh
// must deliver the tunnel process's Wait error: a tunnel that dies before
// becoming ready (e.g. an authentication failure) is reported immediately.
func (t *Tunnel) WaitForReady(timeout time.Duration, exitCh <-chan error) error {
	deadline := time.Now().Add(timeout)

	for {
		if conn, err := net.DialTimeout("tcp", t.LocalAddr(), tunnelDialTimeout); err == nil {
			_ = conn.Close()

			return nil
		}

		select {
		case waitErr := <-exitCh:
			if waitErr != nil {
				return fmt.Errorf("the ssh tunnel exited before becoming ready: %w", waitErr)
			}

			return errors.New("the ssh tunnel exited before becoming ready")
		case <-time.After(tunnelPollInterval):
		}

		if time.Now().After(deadline) {
			return fmt.Errorf("the ssh tunnel did not become ready within %s", timeout)
		}
	}
}

// Stop terminates the tunnel process.
func (t *Tunnel) Stop() error {
	if t.cmd == nil || t.cmd.Process == nil {
		return nil
	}

	return t.cmd.Process.Kill()
}

// freeLocalPort asks the kernel for a free TCP port on the loopback
// interface.
// ponytail: the port is released before ssh binds it, so a race can hand it
// to another process; ExitOnForwardFailure=yes makes ssh fail with a clear
// error instead. Retry-on-bind-failure if it ever matters.
func freeLocalPort() (uint16, error) {
	l, err := net.Listen("tcp", tunnelLocalBindAddr+":0")
	if err != nil {
		return 0, err
	}
	defer l.Close()

	return uint16(l.Addr().(*net.TCPAddr).Port), nil
}
