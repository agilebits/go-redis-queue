package airq

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"io"
	"time"

	"github.com/vmihailenco/msgpack"
)

// Job is the struct of job in queue
type Job struct {
	Body         string    `msgpack:"body"`
	ID           string    `msgpack:"id"`
	Unique       bool      `msgpack:"-"`
	When         time.Time `msgpack:"-"`
	WhenUnixNano int64     `msgpack:"when"`
}

func (j *Job) generateID() string {
	if j.Unique {
		b := make([]byte, 40)
		rand.Read(b)
		return base64.URLEncoding.EncodeToString(b)
	}
	h := sha1.New()
	io.WriteString(h, j.Body)
	return hex.EncodeToString(h.Sum(nil))
}

func (j *Job) setDefaults() {
	if j.ID == "" {
		j.ID = j.generateID()
	}
	if j.When.IsZero() {
		j.When = time.Now()
	}
	j.WhenUnixNano = j.When.UnixNano()
}

func (j *Job) String() string {
	j.setDefaults()
	b, _ := msgpack.Marshal(j)
	return string(b)
}
