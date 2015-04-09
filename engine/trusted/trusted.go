package trusted

import (
	"fmt"
	"strconv"

	"github.com/Sirupsen/logrus"
)

type Event int

const (
	EVENT_NONE   Event = 0
	EVENT_CREATE       = 1 << iota
	EVENT_DELETE
	EVENT_IMAGE
	EVENT_CONTAINER
	EVENT_CONTROL
)

type Credentials struct {
	Uid int
	Pid int
	Gid int
	Lid int
}

func extractVars(vars map[string]string) Credentials {
	var (
		u int
		p int
		g int
		l int
	)

	u, _ = strconv.Atoi(vars["ruid"])
	p, _ = strconv.Atoi(vars["rpid"])
	g, _ = strconv.Atoi(vars["rgid"])
	l, _ = strconv.Atoi(vars["rlid"])

	return Credentials{u, p, g, l}
}

func lookupUid(uid int) string {
	return fmt.Sprintf("User %d", uid)
}

func lookupGid(uid int) string {
	return fmt.Sprintf("Group %d", uid)
}

func Audit(typ Event, vars map[string]string, context interface{}) {
	credentials := extractVars(vars)
	uname := lookupUid(credentials.Uid)
	lname := lookupUid(credentials.Lid)
	gname := lookupGid(credentials.Gid)
	var logString string

	if typ&EVENT_CREATE != 0 {
		if typ&EVENT_IMAGE != 0 {
			logString = fmt.Sprintf("Image Create: %s")
		} else if typ&EVENT_CONTAINER != 0 {
			logString = fmt.Sprintf("Container Create: %s")
		} else {
			logString = "ERROR: %s"
		}
	} else if typ&EVENT_DELETE != 0 {
		logString = fmt.Sprintf("Delete: %s")
	} else if typ&EVENT_CONTROL != 0 { // container only - start, stop, kill
		logString = fmt.Sprintf("Control Event: %s")
	} else {
		logString = "ERROR: %s"
	}

	logrus.Infof(logString, context)
}
