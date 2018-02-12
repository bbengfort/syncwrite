package syncwrite

import "encoding/json"

// Entry is a simple wrapper for log entry to be written to disk.
type Entry struct {
	Index int
	Value []byte
}

// Dump an entry to JSON bytes
func (e *Entry) Dump() ([]byte, error) {
	return json.Marshal(e)
}

// Load an entry from JSON data
func (e *Entry) Load(data []byte) (err error) {
	return json.Unmarshal(data, e)
}
