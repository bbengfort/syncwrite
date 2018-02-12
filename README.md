# Sync Writes

**Throughput study for multiple threads writing to a single synchronized object.**

## Benchmarks

To run the benchmarks:

```
$ make bench
```

You may have to install dependencies using `godep restore` beforehand. Note that the benchmarks may take a while to run. There are two types of benchmarks: go bench benchmarks the amount of time it takes an operation to run, and throughput benchmarks measure the number of writes per-second.

The operation benchmarks are as follows:

```
BenchmarkInMemoryLog-8   	10000000	       201 ns/op
BenchmarkFileLog-8       	  200000	      7253 ns/op
BenchmarkLevelDBLog-8    	  100000	     12104 ns/op
```

The throughput benchmarks are as follows:

| Throughout (ops/sec)    |    Smallest |     Largest |     Average |    Confidence |
|-------------------------|------------:|------------:|------------:|--------------:|
| In-Memory Single Thread | 3716351.687 | 9193643.882 | 7183448.018 | ± 1762739.927 |
| In-Memory 8 Threads     |  502870.185 |  693422.327 |  583877.489 |   ± 57614.640 |
| File Single Thread      |  133294.372 |  162959.403 |  154304.171 |    ± 8339.887 |
| File 8 Threads          |   15355.984 |   17014.469 |   16061.860 |     ± 427.865 |
| LevelDB Single Thread   |   77433.163 |   94672.341 |   87358.272 |    ± 4462.037 |
| LevelDB 8 Threads       |    9033.938 |    9527.512 |    9298.797 |     ± 143.880 |

See details for more information.

## Details

This library contains several append-only logs that synchronize accesses using `sync.RWMutex`. The log interface is as follows:

```go
type Log interface {
	Open(path string) error
	Append(value []byte) error
	Get(index int) (*Entry, error)
	Close() error
}
```

The `Open`, `Append`, and `Close` methods are all protected by a write lock, and the `Get` method is protected with a read lock. The logs implemented so far are:

- `InMemoryLog`: appends to an in-memory slice and does not write to disk.
- `FileLog`: on open, reads entries from file into in-memory slice and reads from it, writes append to both the slice and the file.
- `LevelDBLog`: both writes and reads go to a LevelDB database.

### Benchmarks

The benchmarks are conducted by running `t` threads, each of which run `n` actions and return the amount of time it takes to run those `n` actions. The write throughput action is simply calling `Write` on the log with `"foo"` as the value.

The total number of operations performed is `n*t` and this is divided by the total number of seconds across all go routines to get the number of operations per second.
