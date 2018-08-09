package airq

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"io"
	"io/ioutil"
	"time"

	"github.com/vmihailenco/msgpack"
)

// Job is the struct of job in queue
type Job struct {
	CompressedContent string    `msgpack:"content"`
	Content           string    `msgpack:"-"`
	ID                string    `msgpack:"id"`
	Unique            bool      `msgpack:"-"`
	When              time.Time `msgpack:"-"`
	WhenUnixNano      int64     `msgpack:"when"`
}

func compress(in string) string {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write([]byte(in))
	gz.Flush()
	gz.Close()
	return b.String()
}

func uncompress(in string) string {
	rdata := bytes.NewReader([]byte(in))
	r, _ := gzip.NewReader(rdata)
	s, _ := ioutil.ReadAll(r)
	return string(s)
}

func (j *Job) generateID() string {
	if j.Unique {
		b := make([]byte, 40)
		rand.Read(b)
		return base64.URLEncoding.EncodeToString(b)
	}
	h := sha1.New()
	io.WriteString(h, j.Content)
	return hex.EncodeToString(h.Sum(nil))
}

func (j *Job) setDefaults() {
	j.CompressedContent = compress(j.Content)
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
