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
 (and using Go =P). But is not a replacement for jq with streaming
 capabilities because it focuses on just projecting a few fields from JSON
 documents in a newline oriented fashion, there is no filtering or any advanced
 features (it is mainly a JSON structured log analyzer).
 
# What
 
jtoh will produce a newline for each JSON document found on the list,
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
 data3:data4
 ```

 It is very limited on what it can do and it is not supposed to save the world =P.
 A more hands on example:
 
 ```
 TODO: example using gcloud
 ```
