package domain

import "fmt"

func (j *Job) Start() error {
	if j.Status != Pending && j.Status != Failed {
		return fmt.Errorf("job cannot start from %s", j.Status)
	}
	j.Status = InProgress
	j.Attempts++
	return nil
}
func (j *Job) Succeed() error {
	if j.Status != InProgress {
		return fmt.Errorf("job is not in progress")
	}
	j.Status = Succeeded
	j.LastError = ""
	return nil
}
func (j *Job) Fail(err error) error {
	if j.Status != InProgress {
		return fmt.Errorf("job is not in progress")
	}
	j.Status = Failed
	if err != nil {
		j.LastError = err.Error()
	}
	return nil
}
