package handler

import (
	"slices"
	"strings"
	"testing"

	servertypes "github.com/paddyoneill/slurm-ui/internal/server/types"
)

func TestValidateCreateJobRequestNormalisesValidInput(t *testing.T) {
	partition := "  debug  "
	environment := []string{"  FOO=bar  ", "PATH=/usr/bin"}
	req := servertypes.CreateJobRequest{
		LaunchConfig: servertypes.LaunchConfig{
			Name:                    "  example-job  ",
			CurrentWorkingDirectory: "  /tmp/workdir  ",
			Partition:               &partition,
			Environment:             &environment,
			CpusPerTask:             intPtr(2),
			TasksPerNode:            intPtr(1),
			MemoryPerNode:           intPtr(512),
			TimeLimit:               intPtr(30),
		},
		Script: "  #!/bin/bash\necho hello  ",
	}

	issues := validateCreateJobRequest(&req)
	if len(issues) != 0 {
		t.Fatalf("expected no validation issues, got %v", issues)
	}

	if req.Name != "example-job" {
		t.Fatalf("expected trimmed name, got %q", req.Name)
	}
	if req.CurrentWorkingDirectory != "/tmp/workdir" {
		t.Fatalf("expected trimmed cwd, got %q", req.CurrentWorkingDirectory)
	}
	if req.Partition == nil || *req.Partition != "debug" {
		t.Fatalf("expected trimmed partition, got %#v", req.Partition)
	}
	if req.Script != "#!/bin/bash\necho hello" {
		t.Fatalf("expected trimmed script, got %q", req.Script)
	}
	if req.Environment == nil {
		t.Fatal("expected environment to be normalised")
	}
	if got, want := *req.Environment, []string{"FOO=bar", "PATH=/usr/bin"}; strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("expected normalised environment %v, got %v", want, got)
	}
}

func TestValidateCreateJobRequestRejectsInvalidValues(t *testing.T) {
	partition := "   "
	environment := []string{"INVALID"}
	req := servertypes.CreateJobRequest{
		LaunchConfig: servertypes.LaunchConfig{
			Name:                    "   ",
			CurrentWorkingDirectory: "relative/path",
			Partition:               &partition,
			Environment:             &environment,
			CpusPerTask:             intPtr(0),
			TasksPerNode:            intPtr(0),
			MemoryPerNode:           intPtr(0),
			TimeLimit:               intPtr(0),
		},
		Script: "   ",
	}

	issues := validateCreateJobRequest(&req)
	if len(issues) != 9 {
		t.Fatalf("expected 9 validation issues, got %d: %v", len(issues), issues)
	}

	gotFields := make([]string, 0, len(issues))
	for _, issue := range issues {
		gotFields = append(gotFields, issue.Field)
	}

	wantFields := []string{
		"name",
		"currentWorkingDirectory",
		"partition",
		"cpusPerTask",
		"tasksPerNode",
		"memoryPerNode",
		"timeLimit",
		"environment",
		"script",
	}

	for _, field := range wantFields {
		if !slices.Contains(gotFields, field) {
			t.Fatalf("expected validation issues to include field %q, got %v", field, gotFields)
		}
	}
}

func TestValidateCreateNotebookRequestRequiresBaseEnvAndNormalisesLaunchConfig(t *testing.T) {
	partition := "  gpu  "
	environment := []string{"  FOO=bar  "}
	req := servertypes.CreateNotebookRequest{
		LaunchConfig: servertypes.LaunchConfig{
			Name:                    "  notebook  ",
			CurrentWorkingDirectory: "  /home/user  ",
			Partition:               &partition,
			Environment:             &environment,
		},
		BaseEnv: nil,
	}

	issues := validateCreateNotebookRequest(&req)
	if len(issues) != 1 {
		t.Fatalf("expected exactly one validation issue, got %d: %v", len(issues), issues)
	}
	if issues[0].Field != "baseEnv" {
		t.Fatalf("expected baseEnv issue, got %+v", issues[0])
	}
	if req.Name != "notebook" || req.CurrentWorkingDirectory != "/home/user" {
		t.Fatalf("expected launch config to be normalised, got name=%q cwd=%q", req.Name, req.CurrentWorkingDirectory)
	}
	if req.Partition == nil || *req.Partition != "gpu" {
		t.Fatalf("expected trimmed partition, got %#v", req.Partition)
	}
	if req.Environment == nil || len(*req.Environment) != 1 || (*req.Environment)[0] != "FOO=bar" {
		t.Fatalf("expected normalised environment, got %#v", req.Environment)
	}
}

func intPtr(value int) *int {
	return &value
}
