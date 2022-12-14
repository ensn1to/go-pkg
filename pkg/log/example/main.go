package exapmle

import (
	"context"

	log "git.enn-edge.com/device_manage/public/log.git"
	"github.com/gofrs/uuid"
)

func main() {
	log.Infof("default\n")

	nOpts := &log.Options{
		Level:             log.InfoLevel.String(),
		DisableCaller:     false,
		DisableStacktrace: false,
		Format:            "console",
		EnableColor:       true,
		Development:       false,
		OutputPaths:       []string{"dm.log", "stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}

	log.ResetDefault(nOpts)

	log.Infof("after new opt\n")

	uid, _ := uuid.NewV4()
	ctx := context.WithValue(context.Background(), "requestID", uid.String())

	log.C(ctx).Infof("test requestID\n")
}
