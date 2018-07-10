// Package fasta provides types to read and write FASTA-encoded files.
package fasta

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

// Sequence is the common interface for a sequence that can be represented in
// FASTA encoding.
type Sequence interface {
	Name() string
	Seq() []byte
}

// Record is a concrete implementation of Sequence and corresponds to a token
// in a FASTA encoded file.
type Record struct {
	Header   string
	Sequence []byte
}

// Name returns the record header.
func (rec *Record) Name() string {
	return rec.Header
}

// Seq returns the record sequence.
func (rec *Record) Seq() []byte {
	return rec.Sequence
}

// A Reader reads FASTA encoded sequences.
type Reader struct {
	r   *bufio.Reader
	err error
	rec *Record
}

// NewReader returns a new reader that reads from f.
func NewReader(f io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(f)}
}

// Read returns a FASTA record from r. Read always returns either a non-nil
// record or a non-nil error, but not both. After reaching EOF, subsequent
// calls to Read will return a nil record and io.EOF.
func (r *Reader) Read() (*Record, error) {
	// Keep returning EOF after EOF reached.
	if r.err == io.EOF {
		return nil, io.EOF
	}

	for {
		line, err := r.r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			// If no newline at end of file.
			if len(line) > 0 {
				r.rec.Sequence = append(r.rec.Sequence, line...)
			}
			r.err = io.EOF
			return r.rec, nil
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 { // Skip empty lines.
			continue
		}

		if line[0] != '>' {
			if r.rec == nil { // reached sequence before the first header.
				return nil, errors.New("fasta: format error: sequence before header")
			}
			r.rec.Sequence = append(r.rec.Sequence, line...)
			continue
		}
		temp := r.rec
		r.rec = &Record{
			Header:   string(line[1:]),
			Sequence: make([]byte, 0),
		}
		if temp != nil {
			return temp, nil
		}
	}
}

// A Writer writes sequences in a FASTA format.
type Writer struct {
	w     io.Writer
	width int
}

// NewWriter returns a new FASTA format writer that writes to w.
func NewWriter(w io.Writer, width int) *Writer {
	if width == 0 {
		width = 1
	}
	return &Writer{
		w:     w,
		width: width,
	}
}

// Write writes a single sequence in w. It return the number of bytes written
// and any error.
func (w *Writer) Write(s Sequence) (n int, err error) {
	var (
		_n int
	)

	// Write the header.
	n, err = w.w.Write([]byte(">" + s.Name()))
	if err != nil {
		return n, err
	}

	// Write the sequence (width letters at each line).
	for i := 0; i < len(s.Seq()); i++ {
		if i%w.width == 0 {
			_n, err = w.w.Write([]byte("\n"))
			if n += _n; err != nil {
				return n, err
			}
		}
		_n, err = w.w.Write([]byte{s.Seq()[i]})
		if n += _n; err != nil {
			return n, err
		}
	}
	_n, err = w.w.Write([]byte("\n"))
	if n += _n; err != nil {
		return n, err
	}

	return n, nil
}
