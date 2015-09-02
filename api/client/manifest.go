package client

import (
	"github.com/docker/distribution/digest"
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
	outdigest := cmd.Bool([]string{"-digest", "d"}, false, "Computes manifest digest instead of returning the manifest")
	//digestalg := cmd.Bool([]string{"-algo", "a"}, false, "Selects the digest algorithm")
	cmd.Require(flag.Min, 1)

	cmd.ParseFlags(args, true)

	if serverResp, err := cli.call("GET", "/images/"+cmd.Arg(0)+"/manifest", nil, nil); err == nil {
		defer serverResp.body.Close()
		if *outdigest {
			digester := digest.Canonical.New()
			if _, err := io.Copy(digester.Hash(), serverResp.body); err != nil {
				cli.out.Write([]byte("Error writing response to stdout: " + err.Error() + "\n"))
			} else {
				//digester.Digest().Algorithm()
				cli.out.Write([]byte(digester.Digest().String() + "\n"))
			}
		} else {
			if _, err := io.Copy(cli.out, serverResp.body); err != nil {
				cli.out.Write([]byte("Error writing response to stdout: " + err.Error() + "\n"))
			}
		}
	} else {
		cli.out.Write([]byte(err.Error() + "\n"))
	}

	return
}
