# RateLimiter

RateLimiter is a Go package that implements a token bucket rate limiter to control the rate of events, such as API requests, on a per-client basis identified by IP address. This rate limiter ensures that each client can only make a specified number of requests within a defined interval, helping to prevent abuse and ensure fair usage of resources.

## Features

- **Token Bucket Algorithm**: Smooth and efficient rate limiting using the token bucket algorithm.
- **Per-Client Rate Limiting**: Rate limiting based on client IP addresses.
- **Concurrency Safe**: Safe to use in concurrent environments.

## Installation

To use the RateLimiter in your project, you can get the package via:

```sh
go get github.com/yourusername/ratelimiter
```

# Usage
## Creating a RateLimiter
Create a new rate limiter instance by specifying the maximum number of requests allowed (rate) and the interval duration (interval):

```go
import (
    "time"
    "github.com/yourusername/ratelimiter"
)

rl := ratelimiter.NewRateLimiter(5, time.Second)
```
In this example, the rate limiter allows up to 5 requests per second per IP address.

## Checking if a Request is Allowed
Use the Allow method to check if a request from a specific IP address should be allowed:

```go
clientIP := "192.168.1.1"
if rl.Allow(clientIP) {
    // Process the request
    fmt.Println("Request allowed")
} else {
    // Reject the request
    fmt.Println("Request denied")
}
```

## Example
Here's a complete example demonstrating how to use the RateLimiter:

```go
package main

import (
    "fmt"
    "time"
    "github.com/andrii-kosinskyi/ratelimiter"
)

func main() {
    rl := ratelimiter.NewRateLimiter(5, time.Second)
    clientIP := "192.168.1.1"

    for i := 0; i < 10; i++ {
        if rl.Allow(clientIP) {
            fmt.Println("Request allowed")
        } else {
            fmt.Println("Request denied")
        }
        time.Sleep(200 * time.Millisecond)
    }
}
```

## Testing
The package includes tests to ensure correct behavior. To run the tests, use:

```sh
go test
```

## How It Works
Token Bucket Algorithm
The RateLimiter uses the token bucket algorithm to manage the rate of requests. Each IP address has an associated "bucket" that holds a certain number of tokens (permissions to make requests).

- **Refill**: Tokens are added to the bucket at a specified rate until the bucket is full.
- **Consume**: Each request consumes one token from the bucket.
- **Deny**: If the bucket is empty (no tokens available), the request is denied.

## Implementation Details
1. Structure: The RateLimiter struct contains a map of buckets, each identified by a client IP address.
2. Mutex: A mutex (sync.Mutex) is used to ensure thread-safe access to the buckets.
3. Refill Logic: Tokens are refilled based on the elapsed time since the last refill, ensuring smooth and fair rate limiting.
4. Allow Method:
   - Checks if tokens need to be refilled based on the elapsed time.
   - Refills tokens if necessary.
   - Checks if there are available tokens to allow the request.
   - Returns true if the request is allowed and decrements the token count; otherwise, returns false.

## License
This project is licensed under the MIT License. See the LICENSE file for details.

## Contributing
Contributions are welcome! Please fork the repository and create a pull request with your changes.

## Contact
For questions or issues, please contact kosinskiy.andrey@ukr.net .