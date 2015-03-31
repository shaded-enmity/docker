package trusted

import (
	"fmt"
	"net/http"
	"os"
)

const (
	VERSION          = 0x00010001
	VERSION_STRING   = "0.1.0.1" // <majore>.<minor>.<build>.<release>
	PARAMETER_STRING = "trusted"
	HEADER_UID       = "X-Trusted-Client-Id"
	HEADER_EUID      = "X-Trusted-Client-EId"
	CONTAINER_LABEL  = "system_u:system_r:container_t:s0-s0:c0.c%d"
	IMAGE_LABEL      = "system_u:system_r:image_t:s0-s0:c0.c%d"
	// ----------------------------------------------------------------------
	TL_INVALID = iota
	TL_ENFORCING
	TL_PERMISSIVE
)

type TrustLevel int
type Headers map[string]string

type HdrTuple struct {
	header string
	id     *int
}

var (
	Level = TL_PERMISSIVE
)

func GetTrustLevel() TrustLevel {
	return Level
}

func SetTrustLevel(level TrustLevel) void {
	Level = level
}

// ----------------------------------------------------------------------
//   DecorateRequest
//
//   Executed by a trusted client. Add uid/euid headers to headers passed
//   in the `headers` argument.
//

func DecorateRequest(request *http.Request) error {
	uid := os.Getuid()
	euid := os.Geteuid()

	for _, hdr := range []HdrTuple{{HEADER_UID, &uid}, {HEADER_EUID, &euid}} {
		request.Header[hdr.header] = fmt.Sprintf("%d", *hdr.id)
	}

	return nil
}

// ----------------------------------------------------------------------
//   ExtractHeaders
//
//   Eecuted by the daemon to extract uid/euid of the user that issued
//   the specific command.
//
func ExtractHeaders(headers Headers) (int, int, error) {
	var (
		uid  = -1
		euid = -1
	)

	for _, hdr := range []HdrTuple{{HEADER_UID, &uid}, {HEADER_EUID, &euid}} {
		if uid, exists := headers[hdr.header]; exists {
			*hdr.id = int(uid)
		} else {
			if GetTrustLevel() == TL_ENFORCING {
				return uid, euid, fmt.Errorf("Header %q not found!", hdr.header)
			} else {
				log.Logf("[trusted] Header %q not found!", hdr.header)
			}
		}
	}

	return uid, euid, nil
}

// ----------------------------------------------------------------------
//   LabelProcess
//
//   Creates a string with SELinux label command given `uid` and target
//   `pid`.
//
func LabelProcess(uid, pid int) (string, error) {
	var (
		label = CONTAINER_LABEL
	)
	return label, nil
}

// ----------------------------------------------------------------------
//   LabelFile
//
//   Creates a string with SELinux label command given `uid` and target
//   path.
//
func LabelFile(uid int, path string) (string, error) {
	var (
		label = IMAGE_LABEL
	)
	return label, nil
}

func CheckProcess(uid, pid int) error {
	return nil
}

func CheckFile(uid int, path string) error {
	return nil
}
