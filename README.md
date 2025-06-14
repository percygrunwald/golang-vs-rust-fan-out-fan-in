# Golang and Rust fan-out/fan-in comparison

This repo compares the speed and efficiency of naive, unoptimized fan-out/fan-in architectures in Golang and Rust. Speed refers to the rate at which useful work is done by a limited set of threads/goroutines. Efficiency refers to the amount of CPU time and peak memory utilization required to achieve the task.

`main.go` and `main.rs` both implement the following architecture:

```
            ┌───────────────┐
            │ Random Numbers│
            │   Generator   │
            └───────┬───────┘
                    │
                    ▼
            ┌────────────────────┐
            │   produceChan      │
            │ (Buffered Channel) │
            └─────────┬──────────┘
                      │
                      ▼
            ┌───────────────────────────────────────────────────────┐
            │                                                       │
            │                Consumer Goroutines                    │
            │                                                       │
            │ ┌───────────────┐ ┌───────────────┐ ┌───────────────┐ │
            │ │ Square Numbers│ │ Square Numbers│ │ Square Numbers│ │
            │ │   Calculator  │ │   Calculator  │ │   Calculator  │ │
            │ └───────┬───────┘ └───────┬───────┘ └───────┬───────┘ │
            │         │                 │                 │         │
            └─────────▼─────────────────▼─────────────────▼─────────┘
                      │
                      ▼
            ┌────────────────────┐
            │   squareChan       │
            │ (Buffered Channel) │
            └─────────┬──────────┘
                      │
                      ▼
            ┌─────────────────┐
            │ Fan-In Goroutine│
            │   Sum Calculator│
            └───────┬─────────┘
                    │
                    ▼
            ┌────────────────┐
            │ Final Sum      │
            │   Printed      │
            └────────────────┘
```

A generator thread/goroutine generates (efficiently) some random numbers into a "producer channel". The numbers come from a list of pre-generated random numbers to remove any differences in RNG speed between languages (see below for more details).

A set of fan-out threads/goroutines consumes from the "producer channel" and does some "work" by simply squaring the number, which is then produced into a "squares channel".

A fan-in thread/goroutine consumes from the "squares channel" and accumulates the sum of the squares.

## Run and profile the Golang implementation

Using Mac OS `time`:

```
/usr/bin/time -l go run main.go
```

## Run and profile the Rust implementation

Using Mac OS `time`:

```
/usr/bin/time -l cargo run main.rs
```

## Analysis of results

My dev machine is a 2023 Macbook Pro M3 Max 48GB.

### Golang results

A typical nth run of `50 million` values to `3` consumers:

```
percy@Percys-MBP-2:~/Code/fan-out-fan-in$ /usr/bin/time -l go run main.go -n 50000000 -w 3
Final Sum: 2345000000
Finished main thread.
        4.25 real         6.93 user         5.60 sys
            24887296  maximum resident set size
            ...
```

`~4.3s` total runtime (`~11.6e6 values/sec`)
`~12.6s` total CPU time (`~0.26 CPU seconds/million values`)
`~24.9mb` max RSS

### Rust results

A typical nth run of `50 million` values to `3` consumers:

```
$ /usr/bin/time -l cargo run --bin main -- -n 50000000 -w 3
    Finished `dev` profile [unoptimized + debuginfo] target(s) in 0.00s
     Running `target/debug/main -n 50000000 -w 3`
Producing 50000000 values to 3 consumers...
Final Sum: 1350000000
Finished main thread.
        5.65 real        10.72 user         2.08 sys
             1769472  maximum resident set size
             ...
```

`~5.7s` total runtime (`~8.7e6 values/sec`)
`~12.8s` total CPU time (`~0.26 CPU seconds/million values`)
`~1.8mb` max RSS

### Comparison

| Item                            | Golang | Rust  | Golang:Rust | Rust:Golang |
| ------------------------------- | ------ | ----- | ----------- | ----------- |
| Total runtime (s)               | 4.3    | 5.7   | -25%        | +33%        |
| Values per second               | 11.6e6 | 8.7e6 | +33%        | -25%        |
| Total CPU time (s)              | 12.6   | 12.8  | -2%         | +2%         |
| CPU time per million values (s) | 0.26   | 0.26  | 0           | 0           |
| Max RSS (mb)                    | 24.9   | 1.8   | +1283%      | -93%        |

## Comparing random number generation performance

It turns out that Golang's `rand` and Rust's `rand` vary massively in their speed, which is demonstrated in `random_{go.go, rust.rs}`. Dev machine is a 2023 Macbook Pro M3 Max 48GB.

```
$ go run random_go.go -n 10000000
Go: Generated 10000000 random numbers in 66.210625ms
```
```
$ cargo run --bin random -- -n 10000000
    Finished `dev` profile [unoptimized + debuginfo] target(s) in 0.00s
     Running `target/debug/random -n 10000000`
Rust: Generated 10000000 random numbers in 4.466431084s
```
