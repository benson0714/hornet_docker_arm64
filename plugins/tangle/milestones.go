package tangle

import (
	"runtime"

	"github.com/gohornet/hornet/packages/model/tangle"
	"github.com/gohornet/hornet/packages/workerpool"
	"github.com/gohornet/hornet/plugins/gossip"
)

var (
	checkForMilestoneWorkerCount = runtime.NumCPU()
	checkForMilestoneQueueSize   = 10000
	checkForMilestoneWorkerPool  *workerpool.WorkerPool
)

func checkBundleForMilestone(bundle *tangle.Bundle) {
	isMilestone, err := tangle.CheckIfMilestone(bundle)
	if err != nil {
		log.Infof("Invalid milestone detected! Err: %s", err.Error())
		return
	}
	if !isMilestone {
		return
	}

	processValidMilestone(bundle)
}

func processValidMilestone(bundle *tangle.Bundle) {
	Events.ReceivedNewMilestone.Trigger(bundle)

	// Mark all tx of a valid milestone as requested, so they get stored on eviction
	// Warp sync milestone txs are not requested by default => they would get lost
	for _, tx := range bundle.GetTransactions() {
		tx.SetRequested(true)
	}

	solidMsIndex := tangle.GetSolidMilestoneIndex()
	bundleMsIndex := bundle.GetMilestoneIndex()

	tangle.StoreMilestoneInCache(bundle)

	latestMilestoneIndex := tangle.GetLatestMilestoneIndex()
	if latestMilestoneIndex < bundleMsIndex {
		tangle.SetLatestMilestone(bundle)
		Events.LatestMilestoneChanged.Trigger(bundle)
		milestoneSolidifierWorkerPool.Submit(bundleMsIndex)
	}

	if bundleMsIndex > solidMsIndex {
		log.Infof("Valid milestone detected! Index: %d, Hash: %v", bundleMsIndex, bundle.GetMilestoneHash())

		// Request trunk and branch
		gossip.RequestMilestone(bundle)

	} else {
		pruningIndex := tangle.GetSnapshotInfo().PruningIndex
		if bundleMsIndex < pruningIndex {
			// This should not happen! We didn't request it and it should be filtered because of timestamp
			log.Panicf("Synced too far! Index: %d (%v), PruningIndex: %d", bundleMsIndex, bundle.GetMilestoneHash(), pruningIndex)
		}
	}

}
