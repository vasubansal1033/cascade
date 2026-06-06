# Cascade - Build a Toy LSM Tree From Scratch
by [Anirudh Rowjee](https://rowjee.com), Software Engineer - Storage @ Couchbase

## What is this?

We will go in depth to understand how a Log-Structured Merge Tree works — the storage engine architecture behind databases like RocksDB, LevelDB, Cassandra, and Couchbase. By the end of this workshop, you will have written your own LSM Tree in Go that can handle writes, reads, and deletions, persist data to disk in SSTable format, and run a full tiered compaction pass across levels.

We define **read amplification** as the number of disk IO calls required to serve a single read operation.

The central tension we will explore: **what are you paying for reads when writes are relatively cheap?** Every design decision in an LSM Tree is a response to that question. You will feel this tradeoff directly, through a read IO counter that tracks how many IOs a read requires as data accumulates across levels.

## Who is it for?

Engineers with moderate systems programming experience who want to understand how LSM Trees work internally — and why they are designed the way they are. No prior database internals knowledge is required, but comfort with a systems programming language is expected. Basic comfort with file IO and goroutines is helpful but not a prerequisite.

## Format

~4 hours. Hands-on.

## What will we build?

A working LSM Tree in Go with three on-disk levels — an unsorted L0, and two sorted levels (L1, L2). The engine will support:

- `GET`, `UPSERT`, and `DELETE` (with tombstones)
- Flushing to disk in SSTable format
- 1 KB memtable with memtable swap
- Reads across all three levels (amplification factor of 10: L0 = 10 KB, L1 = 100 KB, L2 = 1 MB)
- Checkpointing and recovery
- Tiered compaction across levels

We will also instrument the engine with a read IO counter, so that the impact of compaction and read path design on read amplification is directly observable.

## Stages

Run the tests for a specific stage with `go test -run TestStageN ./...`

1. **Setup** — Get the skeleton repository running, understand the testing harness, and run the first milestone tests. Sets the feedback loop for the rest of the workshop.
2. **Flush** — Persist data to disk as an SSTable. Introduce the on-disk format and discuss encoding choices.
3. **L0 Reads** — Implement reads from L0 SSTables. Introduce the read IO counter. See directly what a read costs when data is scattered across SSTables in a single level.
4. **Checkpointing and Recovery** — Handle restarts gracefully: snapshot current state to disk with `Sync` and restore it on startup with `Restart`.
5. **Compaction in Isolation** — Understand the compaction algorithm as a standalone problem: merging and sorting two SSTable streams, handling tombstones, producing a new SSTable.
6. **Tiered Compaction End-to-End** — Wire compaction into the engine. Implement a full tiered compaction pass across L0, L1, and L2. Observe the reduction in read IOs as compaction consolidates SSTables.

## Outcome

By the end of this workshop you will understand:

- The write-optimised tradeoff at the heart of LSM Trees, and why it exists
- The on-disk layout of an SSTable and how it is written and read
- Why L0 is unsorted, and what that costs on the read path
- How tombstones work and why deletes are non-trivial
- How read amplification grows as SSTables accumulate, and how compaction recovers it
- The tiered compaction algorithm: what it costs in write amplification, what it saves in implementation complexity
- What real systems like LevelDB and RocksDB extend from this foundation

## Tooling

Go. A skeleton repository with type definitions, method stubs, and per-stage tests is provided — you write the implementations.

Run a specific stage:

```
go test -run TestStage1 ./...   # Setup
go test -run TestStage2 ./...   # Flush
go test -run TestStage3 ./...   # L0 Reads
go test -run TestStage4 ./...   # Checkpointing and Recovery
go test -run TestStage5 ./...   # Compaction in Isolation
go test -run TestStage6 ./...   # Tiered Compaction End-to-End
```

## Further Reading

- **O'Neill et al. (1996)** — [The Log-Structured Merge-Tree](https://www.cs.umb.edu/~poneil/lsmtree.pdf) — the canonical paper that introduced LSM Trees. Read this after the workshop to see how much of what you built was already in the original design.
- **LevelDB source code** — [github.com/google/leveldb](https://github.com/google/leveldb) — a direct, readable implementation of the ideas you just built, and the direct ancestor of RocksDB. Every component will look familiar.
- **RocksDB Wiki: Compaction** — [github.com/facebook/rocksdb/wiki/Compaction](https://github.com/facebook/rocksdb/wiki/Compaction) — how a production system extends tiered and leveled compaction, and the operational tradeoffs between them.
