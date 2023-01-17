package pkg

import (
	"log"
	"sync"
)

type Record struct {
	sync.RWMutex

	content []string
}

func NewRecord(content []string) Record {
	return Record{
		content: []string{},
	}
}

func (r *Record) Append(line string) {
	r.Lock()
	defer r.Unlock()
	r.content = append(r.content, line)
}

func (r *Record) Replace(new []string) {
	r.Lock()
	defer r.Unlock()

	r.content = make([]string, len(new))

	copied := copy(r.content, new)
	log.Printf("Copied %d elements form new content(%d)", copied, len(new))
}

func (r *Record) GetLine() string {
	r.RLock()
	defer r.RUnlock()

	// TODO: Needs to make sure the record is not empty
	ans := r.content[0]
	r.content = r.content[1:]

	log.Println("Get line")
	return ans
}
