package jobs

import (
	"encoding/json"
	"log"

	"github.com/krysopath/klystron/structs"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ValidateJob(job *structs.Job) bool {
	valid := true
	dir := job.Directory
	if string(dir) != dir {
		log.Fatal("job has weird base directory")
		valid = false
	}
	if len(job.Directory) < 1 {
		log.Println("job has weird base directory")
		valid = false
	}
	if len(job.Outputs) < 1 {
		log.Println("job has no specified outputs")
		valid = false
	}
	return valid

}

func JobMarshal(job *structs.Job) []byte {
	bytes, err := json.Marshal(job)
	check(err)
	return bytes
}

func JobUnmarshal(data []byte) structs.Job {
	var job structs.Job
	err := json.Unmarshal(data, &job)
	check(err)
	if validateJob(&job) {
		return job
	} else {
		panic("job did not validate")
	}
}
