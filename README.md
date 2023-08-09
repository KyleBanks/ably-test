# Ably Test

## Dependencies

- Go
- Node.js
   - +Jest: `npm install jest --global`

## Testing

Run both the Go and Node unit tests:

```
$ make test
```

*Given limited time, only one test has been written for each as a demonstration.*

## Running

Before running the client or server application, be sure to set the Ably API key environment variable:

```
$ export ABLY_API_KEY="redacted";
```

Run the server in one terminal:

```
$ make server
```

And then in two or more terminals, running the client application: 

```
$ make client
```

See the [Makefile](./Makefile) for more commands.