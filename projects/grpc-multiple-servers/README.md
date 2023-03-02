<!--forhugo
+++
title="GRPC Client and Multiple Servers Communication"
+++
forhugo-->

Timebox: 5 days

## Learning objectives:

- Writing the proto definitions from scratch
- Hosting multiple services
- Single binaries acting as both clients and servers
- Contacting multiple servers in parallel and handling unsuccessful responses

## What is gRPC?

Read the [gRPC Introduction](https://grpc.io/docs/what-is-grpc/introduction/) and
[gRPC Core Concepts](https://grpc.io/docs/what-is-grpc/core-concepts/) for an overview of gRPC.

## Build a gRPC based prober

```console
> protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    prober/prober.proto
```

Observe the new generated files:

```console
> ls prober
> prober.pb.go		prober.proto		prober_grpc.pb.go
```

Read through `prober_grpc.pb.go` - this is the interface we will use in our code. This is how gRPC works: a `proto` format
gets generated into language-specific code that we can use to interact with gRPCs. If we're working with multiple programming languages
that need to interact through RPCs, we can do this by generating from the same protocol buffer definition with the language-specific
tooling.

Now run the server and client code. You should see output like this.

```console
> go run prober_server/main.go
> 2022/10/19 17:51:32 server listening at [::]:50051
```

```console
> go run prober_client/main.go
> 2022/10/19 17:52:15 Response Time: 117.000000
```

We've now gained some experience with the protocol buffer format, learned
how to generate Go code from protocol buffer definitions, and called that code from Go programs.

## Implement prober logic

Let's modify the prober service slightly. Instead of the simple one-off HTTP GET against a hardcoded google.com, we are going to modify the service to probe an HTTP endpoint N times and return the average time to GET that endpoint to the client.

Change your prober request and response:

- Add a field to the `ProbeRequest` for the number of requests to make.
- Rename the field in `ProbeReply` (and perhaps add a comment) to make clear it's the _average_ response time of the several requests.

Note that it's ok to rename fields in protobuf (unlike when we use JSON), because the binary encoding of protobuf messages doesn't include field names.
However, do note that we need to be very careful about removing proto fields, changing their types, or changing the numerical ordering of fields. In general, these kinds of changes will break your clients because they change the binary encoding format.
You can [read more about backward/forward compatibility with protobufs](https://earthly.dev/blog/backward-and-forward-compatibility/) if you want.

Remember that you'll need to re-generate your Go code after changing your proto definitions.

Update your client to read the endpoint and number of repetitions from the [command line](https://gobyexample.com/command-line-arguments).
Then update your server to execute the probe: do a HTTP fetch of `endpoint` the specified number of times.
The initial version of the code demonstrates how to use the standard [`net/http` package](https://pkg.go.dev/net/http) and the standard time package.

Add up all the elapsed times, divide by the number of repetitions, and return the average to the client.
The client should print out the average value received.
You can do arithmetic operations like addition and division on `time.Duration` values in Go.

## Add a client timeout

Maybe the site we are probing is very slow (which can happen for all kinds of reasons, from network problems to excessive load),
or perhaps the number of repetitions is very high.
Either way, we never want our program to wait forever.
If we are not careful about preventing this then we can end up building systems where problems in one small part of the system
spread across all of the services that talk to that part of the system.

On the client side, add a [timeout](https://pkg.go.dev/context#WithTimeout) to stop waiting after 1 second.

Run your client against some website - how many repetitions do you need to see your client timeout?

## Handling Errors

How do we know if the HTTP fetch succeeded at the server? Add a check to make sure it did.

How should we deal with errors, e.g. if the endpoint isn't found, or says the server is in an error state?
Modify your code and proto format to handle these cases.

## Extra Challenge: Serve and Collect Prometheus Metrics

These sections are optional - do them for an extra challenge if time permits.

### Part 1: Add Prometheus Metrics

Let's learn something about how to monitor applications.

In software operations, we want to know what our software is doing and how it is performing.
One very useful technique is to have our program export metrics. Metrics are basically values that your
program makes available (the industry standard is to export and scrape over HTTP).

Specialised programs, such as Prometheus, can then fetch metrics regularly
from all the running instances of your program, store the history of these metrics, and do useful arithmetic on them
(like computing rates, averages, and maximums). We can use this data to do troubleshooting and to alert if things
go wrong.

Read the [Overview of Prometheus](https://prometheus.io/docs/introduction/overview/).

Now add Prometheus metrics to your prober server. Every time you execute a probe, update a `gauge` metric that tracks the latency.
Add a `label` specifying the endpoint being probed.
The [Prometheus Guide to Instrumenting a Go Application](https://prometheus.io/docs/guides/go-application/) has all the information you need to do this.

Once you've run your program, use your client to execute probes against some endpoint.
Now use the `curl` program or your browser to view the metrics.
You should see a number of built-in Go metrics, plus your new gauge.

If you use your client to start probing a second endpoint, you should see a second labelled metric appear.

### Part 2: Add distributed tracing

When working with complex distributed systems, we always need to keep in mind how we plan to monitor and observe that our system is behaving properly in production.

Requests will not take linear path. In our simple case we can not guarantee that among multiple servers, a particular one will handle this exact request.

Luckily, there is a open-source Observability framework OpenTelemetry to help us. Get familiar with the concept (Observability)[https://opentelemetry.io/docs/concepts/observability-primer/#what-is-observability].

In this section we will add distributed tracing. We will use Honeycomb, a SaaS distributed tracing provider.
Add tracing to all parts of your application.

Run your system and view traces in Honeycomb. Honeycomb provides a useful guide for their own OpenTelemetry Distribution for Go. Run through the Honeycomb sandbox tour and then explore your own data in the same way.

Of course, this application is too simple for using tracing in real life, but have you noticed how simple the process went? Did you need to generate any manual spans? Now have a quick research about tracing RPC vs gGRPC. Discuss it with your mentor or in the channel. The purpose of this exercise is for you to practice two things:

- Practice adding tracing to simple application before we jump into more complex stuff
- Making a cautious choice of frameworks and libraries for your project implementation keeping in mind how difficult or easy it will be to observe the behavior of the service you provide.
