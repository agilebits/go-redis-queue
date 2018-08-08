package airq

import "testing"

func TestSetDefaults(t *testing.T) {
	t.Parallel()
	j := &Job{}
	j.setDefaults()
	if j.ID == "" {
		t.Error("job.ID should be generated")
	}
	if j.When.IsZero() {
		t.Error("job.When should be now")
	}
}
