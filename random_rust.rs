use rand::{Rng, thread_rng};
use std::time::Instant;

fn main() {
    let mut rng = thread_rng();
    let args: Vec<String> = std::env::args().collect();
    let n_arg = args
        .iter()
        .position(|arg| arg == "-n")
        .and_then(|pos| args.get(pos + 1));
    let n: usize = match n_arg {
        Some(val) => val.parse().expect("Please provide a valid number after -n"),
        None => {
            eprintln!("Error: No value provided for -n. Please pass a value using -n <number>");
            std::process::exit(1);
        }
    };

    let start = Instant::now();

    for _ in 0..n {
        let _ = rng.gen_range(1..=10); // Generate random number between 1 and 10
    }

    let duration = start.elapsed();
    println!("Rust: Generated {} random numbers in {:?}", n, duration);
}
