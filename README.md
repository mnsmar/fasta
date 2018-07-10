# fasta

## Installation

```bash
go get github.com/mnsmar/fasta
```

## Description
A library that provides handling of FASTA-encoded files for the Go language.

## Example

```go
f := "foo.fa"
r := NewReader(os.Open(f))

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
```
