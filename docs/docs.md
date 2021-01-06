<details>
<summary>GET /queues</summary>

Get information about every queue,
ordered by order created





##### Responses
- `200`

  A list of queues

  #### Body
```json
{
    "queues": [
        {
            "id": "queue UUID",
            "name": "queue name",
            "capacity": "capacity of this queue",
            "ephemeral": "is this queue ephemeral?",
            "size": "int size of this queue"
        },
        "..."
    ],
    "count": {
        "queues": "int length of the queues list"
    }
}
```



</details>

<details>
<summary>POST /queues</summary>

Create a queue

###### Body
|name|type|description|default|
| - | - | - | - |
|name|optional string|Name of this queue. Can be used instead interchangeably with its id in API calls, and so it should be unique|`null`|
|capacity|optional int|The capacity of this queue If messages are pushed onto a full queue, whatever is on the head is pushed out in a fifo style. If the capacity is null, the queue has an unlimited size|`null`|
|ephemeral|bool|Is this queue ephemeral? Ephemeral queues are not backed by any storage, but instead are completely in memory. This allowes them to be written to read read from quickly, but they are lost when the server goes down|`false`|


#### Body
```json
{
    "name": "queue name",
    "capacity": "capacity of this queue",
    "ephemeral": "is this queue ephemeral?"
}
```


##### Responses
- `200`

  The queue was created

  #### Body
```json
{
    "id": "queue UUID"
}
```


- `409`

  A queue of this name already exists

  #### Body
```json
{
    "error": "conflict"
}
```



</details>

<details>
<summary>GET /queues/:queue</summary>

Get information about this queue





##### Responses
- `200`

  Information about this queue

  #### Body
```json
{
    "queue": {
        "id": "queue UUID",
        "name": "queue name",
        "capacity": "capacity of this queue",
        "ephemeral": "is this queue ephemeral?",
        "size": "int size of this queue"
    }
}
```


- `400`

  No such queue exists

  #### Body
```json
{
    "error": "no_queue"
}
```



</details>

<details>
<summary>PUT /queues/:queue</summary>

Put a message onto the tail of this queue
If this queue does not exist, it will be created
populated with the sent message





##### Responses
- `201`

  The message was put onto the queue's tail

  

- `406`

  A message couldn't be enqueued because the queue is at capacity

  #### Body
```json
{
    "error": "at_capacity"
}
```


- `411`

  The length header is missing

  #### Body
```json
{
    "error": "length_required"
}
```


- `413`

  The message is too long to fit on the queue

  #### Body
```json
{
    "error": "message_too_long"
}
```



</details>

<details>
<summary>DELETE /queues/:queue</summary>

Delete this queue
Multiple calls are idempotent, so if there is no queue targeted
nothing will happen





##### Responses
- `204`

  Queue was deleted

  


</details>

<details>
<summary>GET /queues/:queue/consume/:index</summary>

Consume a message from this queue
Consuming a message will return its content,
and delete it from the queue.
Indexing is 0 based, starting from the head,
where the most recent message will be





##### Responses
- `200`

  Whatever is on the queue at this index

  #### Body
```json
{
    "message": "enqueue'd message"
}
```


- `404`

  No message exists on the queue here

  #### Body
```json
{
    "error": "no_message"
}
```



</details>

<details>
<summary>GET /queues/:queue/head</summary>

Pop the next message on this queue
This will consume the item at the queue's head
This is equivalent to `GET /queues/:queue/consume/0`





##### Responses
- `200`

  Whatever is on this queue's head

  #### Body
```json
{
    "message": "enqueue'd message"
}
```


- `404`

  No message exists on the queue here

  #### Body
```json
{
    "error": "no_message"
}
```



</details>

<details>
<summary>GET /queues/:queue/peek/:index</summary>

Read a message from this queue without consuming it
This functions similarly to `/consume`,
but does not consume messages when they are read.
Thus allowing you to "peek" at messages





##### Responses
- `200`

  Whatever is on the queue at this index

  #### Body
```json
{
    "message": "enqueue'd message"
}
```


- `404`

  No message exists on the queue here

  #### Body
```json
{
    "error": "no_message"
}
```



</details>

<details>
<summary>GET /queues/:queue/tail</summary>

Get the last message on this queue
This will consume the items at the queue's tail
This is equivalent to `GET /queues/:queue/consume/<len>`
where `<len>` == this queue's length - 1





##### Responses
- `200`

  Whatever is on this queue's tail

  #### Body
```json
{
    "message": "enqueue'd message"
}
```


- `404`

  No message exists on the queue here

  #### Body
```json
{
    "error": "no_message"
}
```



</details>
