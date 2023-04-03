# Lazy Reader At
This package contains a function to turn an `io.Reader` into a `io.ReaderAt` lazily.
The traditional approach would be to read the whole stream into a `bytes.Buffer` and using it's `ReadAt` method. The disadvantage of that approach is, that even if the consumer of `io.ReaderAt` just reads the start of the data, it will still need to load the whole stream into memory. The lazy reader at will only load the parts that are actually accessed (and everything before it). So if the consumer just accesses parts at the start or in the middle of the stream, this will save you valuable memory.

## Usage
Using the wrapper is pretty easy. Just call `ra := lra.NewLazyReaderAt(r)` where `r` is your `io.Reader` and `ra` is the resulting `io.ReaderAt`.

```go
// random is a source of infinitely many pseudo-random bytes
random := rand.New(rand.NewSource(42))
readerAt := NewLazyReaderAt(random)

buf := make([]byte, 8)
n, err := readerAt.ReadAt(buf, 100) // skip the first 100 bytes and then read 8 bytes
// handle the error and n

// now readerAt has read 108 bytes into its buffer. We still have access to the previously read bytes

n, err = readerAt.ReadAt(buf, 50) // read 8 bytes from the middle. The data was buffered when we first called ReadAt
// handle the error and n

random2 := rand.New(rand.NewSource(123))
readerAt.Reset(random2) // you can recycle an instance (and it's buffer) by calling the Reset method and passing a new io.Reader
```
