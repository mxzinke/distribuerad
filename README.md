# Distribuerad

> A tool to assign and manage tasks, distributed to multiple services.

**This project is heavily under construction, no actual use yet possible.**

It addresses the issue, when you want to spawn a new task in a distributed system, which should be
executed next (or with delay). The same thing does every Pub/Sub System, the only thing which is different,
is that the events can be delayed.

Distribuerad is design to use very low resources and can easily deploy into a Kubernetes environment 
as a Docker Container.

## Deploy

*This part will come very soon...*

## How to use it?

It's very easy. We only have two main endpoints. The first is to publish/post a new task or event. The second
is to get the list of open tasks/events. By using the third endpoint, the event can be removed.

### Add new events

Using the following endpoint, you can create a new event in a channel. You don't have to handle, about
creating the channel.

```bash
POST /event/:channel-name
HEADER: [ "Content-Type": "application/json" ]
BODY: {
  "data": "<string data (e.g. JSON) which will be given to the service>",
  # Optional Parameters:
  "publishAt": "<In ISO-8601 format, if not given it will publish immediately>"
}
```

### Check the list of events

```bash
GET /event/:channel-name
```

You will receive a list of events.
