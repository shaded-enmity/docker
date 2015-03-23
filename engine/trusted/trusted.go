package trusted;

const (
	VERSION = 0x00010001,
	VERSION_STRING = "0.1.0.1", // <majore>.<minor>.<build>.<release>
	PARAMETER_STRING = "trusted",
	HEADER_UID  = "X-Trusted-Client-Id",
	HEADER_EUID = "X-Trusted-Client-EId"
	// ----------------------------------------------------------------------
	TL_INVALID = iota
	TL_ENFORCING
	TL_PERMISSIVE
);

type TrustLevel int;
type Headers map[string]string;

type HdrTuple struct {
	header string
	id     int*
};

var (
	Level = TL_ENFORCING
);

func GetTrustLevel() TrustLevel {
	return Level;
}

func SetTrustLevel(level TrustLevel) void {
	Level = level;
}

func DecorateRequest(headers Headers) (Headers, error) {
	uid := os.Getuid();
	euid := os.Geteuid();
	
	for _, hdr := range []HdrTuple { {HEADER_UID, &uid}, {HEADER_EUID, &euid} } {
		headers[hdr.header] = fmt.Sprintf("%d", *hdr.id);
	}

	return headers, nil
}

func ExtractHeaders(headers Headers) (int, int, error) {
	var (
		uid  = -1,
		euid = -1
	)

	for _, hdr := range []HdrTuple { {HEADER_UID, &uid}, {HEADER_EUID, &euid} } {
		if uid, exists := headers[hdr.header]; exists {
			*hdr.id = int(uid)
		} else {
			if GetTrustLevel() == TL_ENFORCING {
				return (uid, euid, error.Newf("Header %q not found!", hdr.header))
			} else {
				log.Logf("[trusted] Header %q not found!", )
			}
		}
	}

	return (uid, euid, nil)
}

func LabelProcess(uid, pid int) (string, error) {
	
}

func LabelFile(uid int, path string) (string, error) {

}

func CheckProcess(uid, pid int) error {
	return nil
}

func CheckFile(uid int, path string) error {
	return nil
}
