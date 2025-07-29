# Server with Opentelemetry instrumentation playground

## How does it work
`IDENTITY=<id> go run .`
where `id` is a number, will start a server with that Identity.
If not specified the identity of the server will be zero.

The server listening port will be 3000 + <id>,
so for instance `IDENTITY=1 go run .` will listen on port 3001.

It exposes a `/countdown/<x>` endpoint, where <x> is a number.
When this endpoint is hit, it will:
1. log the received call
2. only if <x> is greater than zero, it will send a call to the next 
   server with identity <id+1> on the endpoint `/countdown/<x-1>`.
3. will return status 200.

So for instance:
```
go run .
IDENTITY=1 go run .
IDENTITY=2 go run .
```
will start 3 servers listening on ports 3000, 3001 and 3002.
And with
`curl -iv http://localhost:3000/countdown/2`
all the 3 servers will be hit in cascade.

`GET http://localhost:3000/countdown/2 -> server id=0 -> GET http://localhost:3001/countdown/1 
-> server id=1 -> GET http://localhost:3002/countdown/0 -> server id=2`

## Instrumentation
Once the request reaches the first server a TraceId will be generated.
The chain of calls to the downstream services will propagate the same TraceId 
and each span will have a different spanId.
Traces will be logged.