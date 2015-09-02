package client

import (
	Cli "github.com/docker/docker/cli"
	flag "github.com/docker/docker/pkg/mflag"
	"io"
)

// CmdVersion shows Docker version information.
//
// Available version information is shown for: client Docker version, client API version, client Go version, client Git commit, client OS/Arch, server Docker version, server API version, server Go version, server Git commit, and server OS/Arch.
//
// Usage: docker version
func (cli *DockerCli) CmdManifest(args ...string) (err error) {
	cmd := Cli.Subcmd("manifest", nil, "Show the Docker version information.", true)
	cmd.Require(flag.Min, 1)

	cmd.ParseFlags(args, true)

	if serverResp, err := cli.call("GET", "/images/"+cmd.Arg(0)+"/manifest", nil, nil); err == nil {
		defer serverResp.body.Close()
		io.Copy(cli.out, serverResp.body)
	} else {
		cli.out.Write([]byte(err.Error() + "\n"))
	}

	return
}
