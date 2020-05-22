package commands

import (
	"strconv"

	"github.com/lefuturiste/jobatator/pkg/store"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var manager *cron.Cron
var managerStarted bool

// DeclareRecurrentJob -
func DeclareRecurrentJob(cmd CmdInterface) {
	// find the queue
	queue, _ := store.FindQueueBySlug(cmd.Parts[1], cmd.User.CurrentGroup, true)
	// we will add the reccurent job to the array of reccurent jobs
	// see if there is similar recurrent job with the same job type in this queue, if already register just say OK
	job, err := store.FindRecurrentJobByType(queue, cmd.Parts[2])
	if err == nil {
		// update the job
		if job.CronExpression != cmd.Parts[3] {
			job.CronExpression = cmd.Parts[3]
			err = job.Update()
			if err != nil {
				log.Error("An error occured while trying to update a recurrent job.", err)
			}
		}

		ReturnString(cmd, "OK#"+strconv.Itoa(job.EntryID))
		return
	}

	// add cron
	if !managerStarted {
		manager = cron.New()
	}
	entryID, err := manager.AddFunc(cmd.Parts[3], func() {
		log.Debug("Started a job (recurrent)")
		publishParts := cmd.Parts
		publishParts[3] = ""
		result, err := PublishUniversal(publishParts, cmd.User)

		if err != nil {
			log.Error("An error occured while trying to publish a recurrent job.", err)
		}
		log.Debug("Recurrent job got job id: ", result)
	})
	if err != nil {
		ReturnError(cmd, "invalid-cron;"+err.Error())
		return
	}
	if !managerStarted {
		manager.Start()
		managerStarted = true
	}

	// create a recurrent job
	job.EntryID = int(entryID)
	job.Type = cmd.Parts[2]
	job.CronExpression = cmd.Parts[3]
	queue.RecurrentJobs = append(queue.RecurrentJobs, job)

	queue.Update() // we want to update the queues without wiping out all the normal jobs in the queue

	ReturnString(cmd, "OK#"+strconv.Itoa(job.EntryID))
}
