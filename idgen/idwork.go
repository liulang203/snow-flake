package idgen

import "os"

var workId, _ = getClientIp()
var IdWork = NewWorker(int64(workId), int64(os.Getpid()))
