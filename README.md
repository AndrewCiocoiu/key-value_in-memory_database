## A thread-safe, in-memory TCP database - written in Go

### Description:
A custom, Redis-style key-value data store built from scratch over raw TCP sockets. Developed as part of my initiative to learn Go by building projects I find interesting, transitioning from writing client-side scripts to building persistent backend servers.

### Features:
- Accepts raw TCP connections so users can interact with the database via terminal tools like `netcat` or `telnet`.
- Supports a custom command protocol: `SET`, `GET`, `DEL`, and `SAVE`.
- 100% thread-safe: Uses read/write mutexes so multiple concurrent users can query or modify the database at the exact same time without crashing the server.
- Persistent state: Can dump the entire in-memory map to a local text file and automatically restore it when the server reboots.

### Usage:
1. Compile binary:
`go build -o mini-db`

2. Run the database server:
`./mini-db`

3. Open a second terminal window and connect to the server using netcat:
`nc localhost 9000`

4. You can now issue commands directly to the database:
> Welcome to Andrew's in-memory DB!
> SET language Go
> Value set!
> GET language
> Go
> DEL language
> language Sucesfully deleted!
> GET language
> Key does not exist in DB!
> SET username Admin
> Value set!
> SAVE
> Database dumped in dump.txt succesfully!

### What I learned:
- Raw TCP Networking - using the `net` package to handle long-lived connection lifecycles instead of simple HTTP requests.
- Mutexes & Thread Safety - using `sync.RWMutex` to lock and unlock memory safely so concurrent goroutines don't cause panic errors.
- Custom Protocol Parsing - reading raw streams of text, cleaning up cross-platform invisible characters (`\r\n`), and executing logic based on the input.
- File I/O - reading from and writing to disk to persist state across server restarts.

### Future Improvements:
Add a `SETEX` command to implement TTL (Time-To-Live) so keys automatically delete themselves after a certain amount of time. Also, maybe add some basic password authentication so anyone on my local network can't just connect and `DEL` my entire database!
