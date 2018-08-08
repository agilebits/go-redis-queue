package airq

import "testing"

func TestCompress(t *testing.T) {
	t.Parallel()
	if out := uncompress(compress("")); out != "" {
		t.Errorf("compression failed %s != \"\"", out)
	}
	if out := uncompress(compress("test")); out != "test" {
		t.Errorf("compression failed %s != \"test\"", out)
	}
}

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
