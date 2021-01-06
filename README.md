# Distribuerad

> A simple and performant tool to distribute events over a message-queue system.

It addresses the issue, when you want to spawn a new task in a distributed system, which should be executed next **or
with delay**. The same thing does every Pub/Sub System, the only thing which is different, is that the events can be
delayed.

Specifically the delay feature is a unique pin-point, which was intentionally made for scheduling a task.

Benchmarks show that the service can handle up to **2 mio / sec.** new events (tested on MacBook i5 / 16GB).

Distribuerad is design to use very low resources and can easily deploy into a Kubernetes environment as a Docker
Container.

## Deploy

```bash
docker pull mxzinke/distribuerad:latest
```

You can use the Docker-Image, which you can find [here](https://hub.docker.com/repository/docker/mxzinke/distribuerad)
or by cloning the repo and building it yourself. You are allowed to use the image as a container in a Kubernetes-Cluster
or Development.

## How to use it?

It's very easy. We only have two main endpoints. The first is to publish/post a new task or event. The second is to get
the list of open tasks/events. By using the third endpoint, the event can be removed / acquired.

### Add new events

Using the following endpoint, you can create a new event in a channel. You don't have to handle, about creating the
channel.

```bash
POST /events/:channel-name
HEADER: [ "Content-Type": "application/json" ]
BODY: {
  "data": "<string data (e.g. JSON) which will be given to the service>",
  # Optional Parameters:
  "publishAt": "<In ISO-8601 format, if not given it will publish immediately>"
}
```

### Check the list of events

```bash
GET /events/:channel-name
```

You will receive a list of events.

### Take up an event

```bash
DELETE /events/:channel-name/:event-id
```

The given event will not be available anymore. The requested service then owns the event. It is the end of the event
lifecycle.