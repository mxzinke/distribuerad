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

## What about performance?

I did a test creating 30 jobs (with each 1 message/second), and some messages simulated by a simple script, creating
additionally 30 messages/second and picking up around 60 messages/second. Everything running in a docker container of
image `mxzinke/distribuerad:latest`.

The result from `docker stats` on MacBook (Intel i5):

```
CONTAINER ID   NAME               CPU %     MEM USAGE / LIMIT     MEM %     NET I/O           BLOCK I/O         PIDS
f8807181c438   distribuerad_1     0.12%     7.91MiB / 3.847GiB    0.12%     86.2kB / 9.87MB   0B / 0B           27
```

To not for get, it does already store around **20k messages** over 2 channels.

## Concept

It's very easy. We separate in our concept into three types of objectives which you can spawn:

* Events - Meaning new messages, directly published to the queue
* Delayed Events - Messages which will be published to the queue at a given point-of-time
* Jobs - Recurring messages, defined by cronJob / cronTab definitions

## The event / message lifecycle

The following steps shows the lifecycle of an element, which you can work with. The lifecycle is the same for 
delayed events and normal events.

### 1. Add an events (Events + Delayed Events)

Using the following endpoint, you can create a new event in a channel. You don't have to handle, about creating the
channel.

```bash
POST /:channel-name/events
HEADER: [ "Content-Type": "application/json" ]
BODY: {
  "data": "<string data (e.g. JSON) which will be shown in the event>",
  # Parameter is used to add a Delayed-Event:
  "publishAt": "<In ISO-8601 format, if not given it will publish immediately>"
  # (OPTIONAL) Time-to-Live (default: 10minutes)
  "ttl": "<after which duration the event should be removed>"
}
```

### 2. Receive the queue

```bash
GET /:channel-name/events
```

You will receive the complete queue, within the channel. If the channel does not exist yet, it will be created.

### 3. Acquire an event / message

To prevent, that any other client is taking up this event, the event will be acquired, just for this service. Also, it
is possible to unlock the event (in case the client could not work on it and want to put it back in the queue).
*The system does not force to lock the event, before removing it from the queue.*

If multiple clients requesting `LOCK` on the same event, the first client trying it, will get the lock. All clients
after, will fail with response status 410 (Gone).

```bash
PUT /:channel-name/events/:event-id
HEADER: [ "Content-Type": "application/json" ]
BODY: {
  # default value is LOCK, only LOCK and UNLOCK possible
  "action": "<LOCK | UNLOCK>"
  # default are 1 minute:
  "ttl": "<duration, tells about how long the lock lasts>"
}
```

### 4. Remove an event / Message

```bash
DELETE /:channel-name/events/:event-id
```

Before doing this, the client normally acquired the event and worked on it. The given event will not be available
anymore. It is the end of the event lifecycle. It is not required to acquire the event.

## (Cron-) Jobs

As already explained in the concept, jobs does allow it to add events within a cron-job, managed by Distribuerad. 
This explicitly make sense, if you have a task to do in you distributed (multi-service) system, which you want to 
make fault-tolerant and prevent executing the same task on every node (in case the services are scaled horizontal).

In following points, there are the possible endpoints for managing jobs explained.

### Add a job

```bash
POST /:channel-name/jobs
HEADER: [ "Content-Type": "application/json" ]
BODY: {
  "name": "<a job name, has to be unique (preventing multiple of the same job)>"
  "data": "<string data (e.g. JSON) which will be shown in the events later on>",
  "cronDef": "<cronJob or cronTab definition, explained below>"
  # (OPTIONAL) Time-to-Live (default: 10minutes)
  "ttl": "<after which duration the event should be removed>"
}
```

The name has to be unique, so that you can later stop the job and don't create multiple jobs in a distributed system.
The data, will be always the same and equals that for the events later on. The `cronDef` parameter is the cronJob or
cronTab definition and can include the timezone and the actual definition.

Following are some **predefined** schedules, but you can use every normal cronJob definition.

```
Entry                  | Description                                | Equivalent To
-----                  | -----------                                | -------------
@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 1 1 *
@monthly               | Run once a month, midnight, first of month | 0 0 1 * *
@weekly                | Run once a week, midnight between Sat/Sun  | 0 0 * * 0
@daily (or @midnight)  | Run once a day, midnight                   | 0 0 * * *
@hourly                | Run once an hour, beginning of hour        | 0 * * * *
```

Also, you can specify **intervals** with the following annotation. For example, "@every 1h30m10s" would indicate a
schedule that activates after 1 hour, 30 minutes, 10 seconds, and then every interval after that.

```
@every <duration>
```

### List all jobs

```bash
GET /:channel-name/jobs
```

You will receive a list of all jobs, within the channel. If the channel does not exist yet, it will be created.

### Remove a job

```bash
DELETE /:channel-name/jobs/:job-name
```

The given job will not be available anymore. The job will not publish events anymore. You can now create a new job with
the name.
