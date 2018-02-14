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

| Throughout (ops/sec)    | Average     | Smallest    | Largest     | Confidence    |
|-------------------------|-------------|-------------|-------------|---------------|
| In-Memory Single Thread | 6903049.264 | 4979794.484 | 8301421.286 | ± 1083618.111 |
| In-Memory 2 Threads     | 7656919.217 | 5554821.085 | 8349907.901 | ± 910246.600  |
| In-Memory 4 Threads     | 8245179.981 | 6547262.066 | 9016449.611 | ± 902530.722  |
| In-Memory 8 Threads     | 760966.574  | 630571.97   | 849180.528  | ± 60596.707   |
| File Single Thread      | 104328.965  | 90748.114   | 117146.058  | ± 8700.082    |
| File 2 Threads          | 109016.35   | 95666.312   | 115915.51   | ± 6765.247    |
| File 4 Threads          | 110962.885  | 102176.54   | 120682.208  | ± 5146.512    |
| File 8 Threads          | 12237.468   | 10935.618   | 13120.329   | ± 614.699     |
| LevelDB Single Thread   | 77131.492   | 59736.012   | 94753.61    | ± 11032.501   |
| LevelDB 2 Threads       | 89536.821   | 85475.76    | 98335.352   | ± 4000.601    |
| LevelDB 24 Threads      | 83864.954   | 70862.801   | 96357.955   | ± 7235.924    |
| LevelDB 8 Threads       | 10387.872   | 9250.536    | 11215.6     | ± 497.232     |

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
