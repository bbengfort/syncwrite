package syncwrite_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bbengfort/syncwrite"
)

var _ = Describe("Write Throughput", func() {

	var log Log
	var err error
	var tmpDir string
	var path string
	var throughput float64

	BeforeEach(func() {
		// Create a temporary file to write the log to.
		tmpDir, err = ioutil.TempDir("", "syncwrite")
		Ω(err).ShouldNot(HaveOccurred())
		path = filepath.Join(tmpDir, "entries.log")
	})

	AfterEach(func() {
		// Clean up log files
		err = os.RemoveAll(tmpDir)
		Ω(err).ShouldNot(HaveOccurred())
	})

	Describe("In-Memory Log", func() {

		BeforeEach(func() {
			log = new(InMemoryLog)
		})

		Measure("single thread throughput", func(b Benchmarker) {
			throughput, err = Benchmark(log, path, 100000, 1)
			Ω(err).ShouldNot(HaveOccurred())
			b.RecordValue("throughput (ops/sec)", throughput)
		}, 15)

		Measure("2 thread throughput", func(b Benchmarker) {
			throughput, err = Benchmark(log, path, 50000, 1)
			Ω(err).ShouldNot(HaveOccurred())
			b.RecordValue("throughput (ops/sec)", throughput)
		}, 15)

		Measure("4 thread throughput", func(b Benchmarker) {
			throughput, err = Benchmark(log, path, 25000, 1)
			Ω(err).ShouldNot(HaveOccurred())
			b.RecordValue("throughput (ops/sec)", throughput)
		}, 15)

		Measure("8 thread throughput", func(b Benchmarker) {
			throughput, err = Benchmark(log, path, 12500, 8)
			Ω(err).ShouldNot(HaveOccurred())
			b.RecordValue("throughput (ops/sec)", throughput)
		}, 15)

		Measure("16 thread throughput", func(b Benchmarker) {
			throughput, err = Benchmark(log, path, 6250, 16)
			Ω(err).ShouldNot(HaveOccurred())
			b.RecordValue("throughput (ops/sec)", throughput)
		}, 15)

	})

	Describe("File Log", func() {

		BeforeEach(func() {
			log = new(FileLog)
		})

		Measure("single thread throughput", func(b Benchmarker) {
			throughput, err = Benchmark(log, path, 100000, 1)
			Ω(err).ShouldNot(HaveOccurred())
			b.RecordValue("throughput (ops/sec)", throughput)
		}, 15)

		Measure("2 thread throughput", func(b Benchmarker) {
			throughput, err = Benchmark(log, path, 50000, 1)
			Ω(err).ShouldNot(HaveOccurred())
			b.RecordValue("throughput (ops/sec)", throughput)
		}, 15)

		Measure("4 thread throughput", func(b Benchmarker) {
			throughput, err = Benchmark(log, path, 25000, 1)
			Ω(err).ShouldNot(HaveOccurred())
			b.RecordValue("throughput (ops/sec)", throughput)
		}, 15)

		Measure("8 thread throughput", func(b Benchmarker) {
			throughput, err = Benchmark(log, path, 12500, 8)
			Ω(err).ShouldNot(HaveOccurred())
			b.RecordValue("throughput (ops/sec)", throughput)
		}, 15)

		Measure("16 thread throughput", func(b Benchmarker) {
			throughput, err = Benchmark(log, path, 6250, 16)
			Ω(err).ShouldNot(HaveOccurred())
			b.RecordValue("throughput (ops/sec)", throughput)
		}, 15)

	})

	Describe("LevelDB Log", func() {

		BeforeEach(func() {
			log = new(LevelDBLog)
		})

		Measure("single thread throughput", func(b Benchmarker) {
			throughput, err = Benchmark(log, path, 100000, 1)
			Ω(err).ShouldNot(HaveOccurred())
			b.RecordValue("throughput (ops/sec)", throughput)
		}, 15)

		Measure("2 thread throughput", func(b Benchmarker) {
			throughput, err = Benchmark(log, path, 50000, 1)
			Ω(err).ShouldNot(HaveOccurred())
			b.RecordValue("throughput (ops/sec)", throughput)
		}, 15)

		Measure("4 thread throughput", func(b Benchmarker) {
			throughput, err = Benchmark(log, path, 25000, 1)
			Ω(err).ShouldNot(HaveOccurred())
			b.RecordValue("throughput (ops/sec)", throughput)
		}, 15)

		Measure("8 thread throughput", func(b Benchmarker) {
			throughput, err = Benchmark(log, path, 12500, 8)
			Ω(err).ShouldNot(HaveOccurred())
			b.RecordValue("throughput (ops/sec)", throughput)
		}, 15)

		Measure("16 thread throughput", func(b Benchmarker) {
			throughput, err = Benchmark(log, path, 6250, 16)
			Ω(err).ShouldNot(HaveOccurred())
			b.RecordValue("throughput (ops/sec)", throughput)
		}, 15)

	})

})

//===========================================================================
// Per-Operation Benchmarks (with go test bench)
//===========================================================================

func BenchmarkInMemoryLog(b *testing.B) {
	log := new(InMemoryLog)
	log.Open("")
	defer log.Close()

	for n := 0; n < b.N; n++ {
		log.Append([]byte("foo"))
	}
}

func BenchmarkFileLog(b *testing.B) {
	tmpDir, err := ioutil.TempDir("", "syncwrites")
	if err != nil {
		b.Error(err)
	}
	defer os.RemoveAll(tmpDir)

	log := new(FileLog)
	if err := log.Open(filepath.Join(tmpDir, "entries.log")); err != nil {
		b.Error(err)
		return
	}
	defer log.Close()

	for n := 0; n < b.N; n++ {
		log.Append([]byte("foo"))
	}
}

func BenchmarkLevelDBLog(b *testing.B) {
	tmpDir, err := ioutil.TempDir("", "syncwrites")
	if err != nil {
		b.Error(err)
	}
	defer os.RemoveAll(tmpDir)

	log := new(LevelDBLog)
	if err := log.Open(tmpDir); err != nil {
		b.Error(err)
		return
	}
	defer log.Close()

	for n := 0; n < b.N; n++ {
		log.Append([]byte("foo"))
	}
}
