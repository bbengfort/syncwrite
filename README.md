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

| Throughout (ops/sec)    |     Average |    Smallest |     Largest |   Confidence |
|-------------------------|------------:|------------:|------------:|-------------:|
| In-Memory Single Thread | 5253135.838 | 4359456.548 | 5591019.392 | ± 353865.702 |
| In-Memory 2 Threads     | 6326110.013 | 4587467.827 | 7130902.694 | ± 750315.325 |
| In-Memory 4 Threads     | 5358033.001 | 4709761.867 | 5770931.862 | ± 313802.803 |
| In-Memory 8 Threads     | 3630376.542 | 2868603.666 | 4322431.426 | ± 362744.642 |
| In-Memory 16 Threads    | 3495941.168 | 3030034.735 | 4234283.768 | ± 270080.857 |
| File Single Thread      |  109835.346 |  107288.465 |  113584.472 |   ± 1841.514 |
| File 2 Threads          |  106918.977 |   97990.257 |  112586.096 |   ± 3771.965 |
| File 4 Threads          |   120557.93 |   105184.91 |  146225.469 |  ± 13477.076 |
| File 8 Threads          |  113074.283 |  109141.601 |  117840.127 |   ± 2544.440 |
| File 16 Threads         |  109542.505 |  105008.955 |    116745.5 |   ± 2983.509 |
| LevelDB Single Thread   |   82206.692 |   79553.087 |   83783.706 |   ± 1096.305 |
| LevelDB 2 Threads       |   82593.469 |   78239.013 |   86354.997 |   ± 1866.427 |
| LevelDB 4 Threads       |   83295.175 |   80616.994 |   85584.984 |   ± 1459.148 |
| LevelDB 8 Threads       |   70046.095 |   66340.541 |   71619.628 |   ± 1236.291 |
| LevelDB 16 Threads      |   68290.992 |   61172.903 |   71260.022 |   ± 2363.267 |

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

**Update** &mdash; the total number of actions `n` has been fixed per measurement, e.g. if a single thread measurement runs 10,000 actions, then the 4 thread measurement will run 2,500 actions per thread.
