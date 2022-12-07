package gracefulshutdown

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	daemon "github.com/iotaledger/hive.go/daemon/ordered"
	"github.com/gohornet/hornet/packages/logger"
	"github.com/gohornet/hornet/packages/node"
)

// maximum amount of time to wait for background processes to terminate. After that the process is killed.
const WAIT_TO_KILL_TIME_IN_SECONDS = 120

var log = logger.NewLogger("Graceful Shutdown")

var PLUGIN = node.NewPlugin("Graceful Shutdown", node.Enabled, func(plugin *node.Plugin) {
	gracefulStop := make(chan os.Signal)

	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		<-gracefulStop

		log.Warningf("Received shutdown request - waiting (max %d seconds) to finish processing ...", WAIT_TO_KILL_TIME_IN_SECONDS)

		go func() {
			start := time.Now()
			for x := range time.Tick(1 * time.Second) {
				secondsSinceStart := x.Sub(start).Seconds()

				if secondsSinceStart <= WAIT_TO_KILL_TIME_IN_SECONDS {
					processList := ""
					runningBackgroundWorkers := daemon.GetRunningBackgroundWorkers()
					if len(runningBackgroundWorkers) >= 1 {
						processList = "(" + strings.Join(runningBackgroundWorkers, ", ") + ") "
					}

					log.Warningf("Received shutdown request - waiting (max %d seconds) to finish processing %s...", WAIT_TO_KILL_TIME_IN_SECONDS-int(secondsSinceStart), processList)
				} else {
					log.Fatal("Background processes did not terminate in time! Forcing shutdown ...")
				}
			}
		}()

		daemon.ShutdownAndWait()
	}()
})
