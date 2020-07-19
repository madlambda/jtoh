# jtoh

jtoh stands for JSON To Human, basically makes it easier to analyze long
streams of JSON objects.
The main use case is to analyze structured logs from Kubernetes and GCP
stack driver. But it will work with any long list/stream of JSON objects.

# Why ?

There is some good tools to parse JSON, like
[jq](https://stedolan.github.io/jq/manual), which I usually use.
But my problem involved processing long lists of JSON documents,
like this (but much bigger):

```json
[
            {"Name": "Ed", "Text": "Knock knock."},
            {"Name": "Sam", "Text": "Who's there?"},
            {"Name": "Ed", "Text": "Go fmt."},
            {"Name": "Sam", "Text": "Go fmt who?"},
            {"Name": "Ed", "Text": "Go fmt yourself!"}
    ]
```

And jq by default does no stream processing, and the stream mode is not
exactly what I want as can be seen on the
[docs](https://stedolan.github.io/jq/manual/#Streaming) and on this
[post](https://devblog.songkick.com/parsing-ginormous-json-files-via-streaming-be6561ea8671).
To be honest I can't even understand the documentation on how jq streaming
works, so even if it is useful for some scenarios it is beyond me to
understand it properly (and what I read on the blog post does
not sound like fun).

The behavior that I wanted is the exact same behavior as
Go's [json.Decoder.Decode](https://golang.org/pkg/encoding/json/#Decoder.Decode),
which is to handle JSON lists as an incremental decoding of each JSON document
inside the list, done in a streaming fashion, hence this tool was built
(and using Go =P). But it is NOT a replacement for jq with streaming
capabilities because it focuses on just projecting a few fields from JSON
documents in a newline oriented fashion, there is no filtering or any advanced
features and it probably won't handle well complex scenarios, it is meant
for long lists of JSON objects or long streams of JSON objects.

# Install

To install it you will need Go >= 1.13. You can clone the repository and run:

```
make install
```

Or you can just run:

```
go install -i github.com/katcipis/jtoh/cmd/jtoh
```

# What
 
jtoh will produce a newline for each JSON document found on the list/stream,
accepting a projection string as a parameter indicating which fields are going
to be used to compose each newline and what is the separator between each field:
 
```
<source of JSON list> | jtoh "<sep>field1<sep>field2<sep>field3.name"
```

Where **<sep>** is the first character and will be considered the separator,
it is used to separate different field definitions and will also be used
as the separator on the output, this:

```
<source of JSON list> | jtoh ":field1:field2"
```

Will generate an stream of outputs like:

```
data1:data2
data1:data2
```

A more hands on example, lets say you are getting the logs for a specific
application on GCP like this:

```sh
gcloud logging read --format=json --project <your project> "severity>=WARNING AND resource.labels.container_name=myapp"
```

You will probably have a long list of something like this:

```json
{
    "insertId": "h3wh26neb0mcbkeou",
    "labels": {
      "k8s-pod/app": "myapp",
      "k8s-pod/pod-template-hash": "56d4fdf46d"
    },
    "logName": "projects/a2b-exp/logs/stderr",
    "receiveTimestamp": "2020-07-14T13:18:40.681669783Z",
    "resource": {
      "labels": {
        "cluster_name": "k8s-cluster",
        "container_name": "myapp",
        "location": "europe-west3-a",
        "namespace_name": "default",
        "pod_name": "kraken-56d4fdf46d-f9trn",
        "project_id": "someproject"
      },
      "type": "k8s_container"
    },
    "severity": "ERROR",
    "textPayload": "cool log message",
    "timestamp": "2020-07-14T13:18:38.741851348Z"
}
```

In this case the application does no JSON structured logging
(which is perfectly fine in some scenarios),
but there is a lot of data around the
actual application log that can be useful for filtering but after
being used for filtering is pure cognitive noise.

Using jtoh like this:

```
gcloud logging read --format=json --project <your project> "severity>=WARNING AND resource.labels.container_name=myapp" | jtoh :timestamp:textPayload
```

You now get a stream of lines like this: 

```
2020-07-14T13:18:38.741851348Z:cool log message
```

The exact same thing is possible with the stream of JSON objects you
get when the application structure the log entries as JSON and you get the
logs directly from Kubernetes using kubectl like this:

```
TODO
```
