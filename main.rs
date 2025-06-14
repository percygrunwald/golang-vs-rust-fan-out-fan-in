use rand::{Rng, thread_rng};
use std::sync::{Arc, Mutex, mpsc};
use std::thread;

fn main() {
    /*
        Fan-Out and Fan-In Diagram:

        Producer Thread
            ┌───────────────┐
            │ Random Numbers│
            │   Generator   │
            └───────┬───────┘
                    │
                    ▼
            ┌────────────────────┐
            │   produce_tx       │
            │ (Sender Channel)   │
            └─────────┬──────────┘
                      │
                      ▼
            ┌───────────────────────────────────────────────────────┐
            │                                                       │
            │                Consumer Threads                       │
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
            │   square_tx        │
            │ (Sender Channel)   │
            └─────────┬──────────┘
                      │
                      ▼
            ┌─────────────────┐
            │ Fan-In Thread   │
            │   Sum Calculator│
            └───────┬─────────┘
                    │
                    ▼
            ┌────────────────┐
            │ Final Sum      │
            │   Printed      │
            └────────────────┘
    */
    const BUFFER_SIZE: usize = 1_000;
    let args: Vec<String> = std::env::args().collect();
    let n_arg = args
        .iter()
        .position(|arg| arg == "-n")
        .and_then(|pos| args.get(pos + 1));
    let num_values: usize = match n_arg {
        Some(val) => val.parse().expect("Please provide a valid number after -n"),
        None => {
            eprintln!("Error: No value provided for -n. Please pass a value using -n <number>");
            std::process::exit(1);
        }
    };
    let w_arg = args
        .iter()
        .position(|arg| arg == "-w")
        .and_then(|pos| args.get(pos + 1));
    let num_consumers: usize = match w_arg {
        Some(val) => val.parse().expect("Please provide a valid number after -w"),
        None => {
            eprintln!("Error: No value provided for -w. Please pass a value using -w <number>");
            std::process::exit(1);
        }
    };

    // Channels
    let (produce_tx, produce_rx) = mpsc::sync_channel(BUFFER_SIZE);
    let (square_tx, square_rx) = mpsc::sync_channel(BUFFER_SIZE);

    // Predefined list of random numbers generated using rng.gen_range
    let random_numbers: Vec<usize> = (0..10).map(|_| thread_rng().gen_range(1..=10)).collect();
    let list_length = random_numbers.len();

    // Function to cycle through the list
    let get_number = move |index: usize| -> usize { random_numbers[index % list_length] };

    println!(
        "Producing {} values to {} consumers...",
        num_values, num_consumers
    );

    // Producer thread
    thread::spawn(move || {
        for i in 0..num_values {
            let num = get_number(i); // Get number from the list
            produce_tx.send(num).unwrap();
        }
    });

    // Wrap produce_rx in Arc<Mutex>
    let produce_rx = Arc::new(Mutex::new(produce_rx));

    // Consumer threads
    let mut consumer_handles = Vec::new();
    for _ in 0..num_consumers {
        let produce_rx = Arc::clone(&produce_rx); // Clone the Arc to share the Receiver
        let square_tx = square_tx.clone();
        consumer_handles.push(thread::spawn(move || {
            for num in produce_rx.lock().unwrap().iter() {
                // Lock and iterate over the Receiver
                square_tx.send(num * num).unwrap();
            }
        }));
    }

    // Fan-in thread
    let sum = Arc::new(Mutex::new(0));
    let sum_clone = Arc::clone(&sum);
    let fan_in_handle = thread::spawn(move || {
        for square in square_rx {
            let mut total = sum_clone.lock().unwrap();
            *total += square;
        }
    });

    // Wait for consumers to finish
    for handle in consumer_handles {
        handle.join().unwrap();
    }
    drop(square_tx); // Close square_tx to signal fan-in thread

    // Wait for fan-in thread to finish
    fan_in_handle.join().unwrap();

    // Print the final sum
    let final_sum = *sum.lock().unwrap();
    println!("Final Sum: {}", final_sum);
    println!("Finished main thread.");
}
