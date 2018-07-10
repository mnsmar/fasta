package fasta

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"
)

// Test Read
var readTests = []struct {
	Test    string
	Data    string
	Err     string
	Headers []string
	Seqs    []string
}{
	{
		Test: "1-seq",
		Data: "" +
			">Seq1\n" +
			"AAA\n" +
			"BBB\n",
		Headers: []string{"Seq1"},
		Seqs:    []string{"AAABBB"},
	},
	{
		Test: "2-seq",
		Data: "" +
			">Seq1\n" +
			"AAA\n" +
			"BBB\n" +
			">Seq2\n" +
			"CCC\n" +
			"DDD\n",
		Headers: []string{"Seq1", "Seq2"},
		Seqs:    []string{"AAABBB", "CCCDDD"},
	},
	{
		Test: "no newline",
		Data: "" +
			">Seq1\n" +
			"AAA\n" +
			"BBB\n" +
			">Seq2\n" +
			"CCC\n" +
			"DDD",
		Headers: []string{"Seq1", "Seq2"},
		Seqs:    []string{"AAABBB", "CCCDDD"},
	},
	{
		Test: "format error",
		Data: "" +
			"AAA\n" +
			">Seq1\n" +
			"BBB\n",
		Err:     "fasta: format error: sequence before header",
		Headers: []string{""},
		Seqs:    []string{""},
	},
}

func TestRead(t *testing.T) {
	for _, tt := range readTests {
		r := NewReader(strings.NewReader(tt.Data))

		for recIdx := range tt.Headers {
			rec, err := r.Read()

			if tt.Err != "" {
				if err == nil || !strings.Contains(err.Error(), tt.Err) {
					t.Errorf("%s: error %q, want error %q", tt.Test, err.Error(), tt.Err)
				}
				continue
			} else if err != nil {
				t.Errorf("%s: unexpected error %q", tt.Test, err.Error())
				continue
			}

			if rec == nil {
				t.Fatalf("%s: unexpected nil record", tt.Test)
			}

			if rec.Name() != tt.Headers[recIdx] {
				t.Errorf("%s: header=%q want %q", tt.Test, rec.Name(), tt.Headers[recIdx])
			}

			if string(rec.Seq()) != tt.Seqs[recIdx] {
				t.Errorf("%s: seq=%q want %q", tt.Test, string(rec.Seq()), tt.Seqs[recIdx])
			}
		}
	}
}

// Test Write
var writeTests = []struct {
	Test    string
	Records []*Record
	Output  string
	Width   int
}{
	{
		Test: "2-seqs write",
		Records: []*Record{
			&Record{Header: "Seq1", Sequence: []byte("AAABBB")},
			&Record{Header: "Seq2", Sequence: []byte("CCCDDD")},
		},
		Output: ">Seq1\nAA\nAB\nBB\n>Seq2\nCC\nCD\nDD\n",
		Width:  2,
	},
	{
		Test: "0-width write",
		Records: []*Record{
			&Record{Header: "Seq1", Sequence: []byte("AAABBB")},
			&Record{Header: "Seq2", Sequence: []byte("CCCDDD")},
		},
		Output: ">Seq1\nA\nA\nA\nB\nB\nB\n>Seq2\nC\nC\nC\nD\nD\nD\n",
		Width:  0,
	},
}

func TestWrite(t *testing.T) {
	for _, tt := range writeTests {
		b := &bytes.Buffer{}
		w := NewWriter(b, tt.Width)

		for _, rec := range tt.Records {
			_, err := w.Write(rec)
			if err != nil {
				t.Errorf("%s: unexpected error %q", tt.Test, err.Error())
				return
			}
		}

		out := b.String()
		if out != tt.Output {
			t.Errorf("%s: out=%q want %q", tt.Test, out, tt.Output)
		}
	}
}

func ExampleReader() {
	in := ">Seq1\nAAA\nBBB\n"
	r := NewReader(strings.NewReader(in))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", record.Name())
		fmt.Printf("%s\n", record.Seq())
	}
	// Output:
	// Seq1
	// AAABBB
}

func ExampleWriter() {
	b := &bytes.Buffer{}
	r := &Record{
		Header:   "Seq1",
		Sequence: []byte("AAABBB"),
	}

	w := NewWriter(b, 4)

	_, err := w.Write(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(b.String())
	// Output:
	// >Seq1
	// AAAB
	// BB
}
