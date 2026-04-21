package slurm

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	slurmgen "github.com/SlinkyProject/slurm-client/api/v0043"
)

const defaultPath = "PATH=/bin:/usr/bin:/usr/local/bin"

type JobInfo struct {
	JobID int
	Name  *string
	State string
	Host  *string
}

type SubmitJobRequest struct {
	Name                    string
	Script                  string
	Partition               *string
	CurrentWorkingDirectory string
	Environment             *[]string
	CpusPerTask             *int
	TasksPerNode            *int
	MemoryPerNode           *int
	TimeLimit               *int
}

func CancelJob(ctx context.Context, jobID int) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	response, err := client.SlurmV0043DeleteJobWithResponse(ctx, strconv.Itoa(jobID), nil)
	if err != nil {
		return fmt.Errorf("delete job %d: %w", jobID, err)
	}

	if response.JSON200 == nil {
		return unexpectedStatus(fmt.Sprintf("delete job %d", jobID), response.StatusCode(), response.Body)
	}

	return nil
}

func GetJob(ctx context.Context, jobID int) (JobInfo, error) {
	client, err := newClient()
	if err != nil {
		return JobInfo{}, err
	}

	response, err := client.SlurmV0043GetJobWithResponse(ctx, strconv.Itoa(jobID), nil)
	if err != nil {
		return JobInfo{}, fmt.Errorf("get job %d: %w", jobID, err)
	}

	if response.JSON200 == nil {
		return JobInfo{}, unexpectedStatus(fmt.Sprintf("get job %d", jobID), response.StatusCode(), response.Body)
	}

	if len(response.JSON200.Jobs) == 0 {
		return JobInfo{}, fmt.Errorf("get job %d: missing job data", jobID)
	}

	job := response.JSON200.Jobs[0]
	state := "UNKNOWN"
	if job.JobState != nil && len(*job.JobState) > 0 {
		state = string((*job.JobState)[0])
	}

	return JobInfo{
		JobID: jobID,
		Name:  job.Name,
		State: state,
		Host:  job.Nodes,
	}, nil
}

func SubmitJob(ctx context.Context, req SubmitJobRequest) (int, error) {
	client, err := newClient()
	if err != nil {
		return 0, err
	}

	response, err := client.SlurmV0043PostJobSubmitWithResponse(ctx, slurmgen.V0043JobSubmitReq{
		Script: &req.Script,
		Job: &slurmgen.V0043JobDescMsg{
			Name:                    &req.Name,
			Partition:               req.Partition,
			CurrentWorkingDirectory: &req.CurrentWorkingDirectory,
			Environment:             buildEnvironment(req.Environment),
			CpusPerTask:             int32Ptr(req.CpusPerTask),
			TasksPerNode:            int32Ptr(req.TasksPerNode),
			MemoryPerNode:           uint64NoValPtr(req.MemoryPerNode),
			TimeLimit:               uint32NoValPtr(req.TimeLimit),
		},
	})

	if err != nil {
		return 0, fmt.Errorf("submit job: %w", err)
	}

	if response.JSON200 == nil {
		return 0, unexpectedStatus("submit job", response.StatusCode(), response.Body)
	}

	if response.JSON200.JobId == nil {
		return 0, fmt.Errorf("submit job: missing job_id in response")
	}

	return int(*response.JSON200.JobId), nil
}

func buildEnvironment(environment *[]string) *slurmgen.V0043StringArray {
	env := make(slurmgen.V0043StringArray, 0, 1)
	if environment != nil {
		env = append(env, (*environment)...)
	}

	hasPath := false
	for _, value := range env {
		if strings.HasPrefix(value, "PATH=") {
			hasPath = true
			break
		}
	}

	if !hasPath {
		env = append(slurmgen.V0043StringArray{defaultPath}, env...)
	}

	return &env
}

func int32Ptr(value *int) *int32 {
	if value == nil {
		return nil
	}

	result := int32(*value)
	return &result
}

func uint32NoValPtr(value *int) *slurmgen.V0043Uint32NoValStruct {
	if value == nil {
		return nil
	}

	result := int32(*value)
	set := true
	return &slurmgen.V0043Uint32NoValStruct{
		Number: &result,
		Set:    &set,
	}
}

func uint64NoValPtr(value *int) *slurmgen.V0043Uint64NoValStruct {
	if value == nil {
		return nil
	}

	result := int64(*value)
	set := true
	return &slurmgen.V0043Uint64NoValStruct{
		Number: &result,
		Set:    &set,
	}
}
