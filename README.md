### Registration
```
curl -X POST \
  http://127.0.0.1:3000/users \
  -H 'content-type: application/json' \
  -d '{
	"username":"user1",
	"email":"user1@gmail.com",
	"password":"0000"
}'

{
    "id": "470dab2c-b491-46e0-8b43-eb8ebc9f2d7b"
}
```
### Login
```
curl -X POST \
  http://127.0.0.1:3000/login \
  -H 'content-type: application/json' \
  -d '{
	"username":"user1",
	"email":"user1@gmail.com",
	"password":"0000"
}'

{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzM5OTQ5MzksInVzZXJfaWQiOiI0NzBkYWIyYy1iNDkxLTQ2ZTAtOGI0My1lYjhlYmM5ZjJkN2IifQ.hE5NF5Dscdtn-Z2pdIyMBrXVeSBG69xeArr8DqoceWo"
}
```

### Session resource
```
curl -X POST \
  http://127.0.0.1:3000/sessions \
  -H 'authorization: bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzM5OTQ5MzksInVzZXJfaWQiOiI0NzBkYWIyYy1iNDkxLTQ2ZTAtOGI0My1lYjhlYmM5ZjJkN2IifQ.hE5NF5Dscdtn-Z2pdIyMBrXVeSBG69xeArr8DqoceWo' \
  -H 'content-type: application/json' \
  -d '{
	"device_id":"0000",
	"platform":"web",
	"model":"chrome",
	"build":1,
	"name":"ievgen-pc"
}'

{
    "created_at": 1532785517921,
    "id": "2c57526a-3f33-41db-9311-719074024a40",
    "messaging_url": "ws://127.0.0.1:3000/d1d0c2bf-66bd-4376-973b-aa0479fab043"
}
```

### Connect to server 
- connection to server for user 1
- sent rpc method for get all events on connection
```javascript
var socket = new WebSocket("ws://127.0.0.1:3000/d1d0c2bf-66bd-4376-973b-aa0479fab043");
socket.onopen = function() {
   socket.send(JSON.stringify({"method":20, "timestamp":new Date().getTime()}));
};
socket.onmessage = function(event) {
  console.log(event.data);
};
```

#### Create group
- In property user_ids put user ID which will be assign to group
```
curl -X POST \
  http://127.0.0.1:3000/groups \
  -H 'authorization: bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzM5OTQ5MzksInVzZXJfaWQiOiI0NzBkYWIyYy1iNDkxLTQ2ZTAtOGI0My1lYjhlYmM5ZjJkN2IifQ.hE5NF5Dscdtn-Z2pdIyMBrXVeSBG69xeArr8DqoceWo' \
  -H 'content-type: application/json' \
  -d '{
	"name":"group_1",
	"user_ids":["6a468c77-3974-4ef0-a45c-a575c57c9b1b"]
}'

{
    "id": "95ba5e74-810d-49b4-943f-0b2753d688c9"
}
```


#### Received group event
```json
{
    "type":70,
    "timestamp":1532786538556,
    "body":{
      "group_id":"95ba5e74-810d-49b4-943f-0b2753d688c9",
      "name":"group_1",
      "user_ids":["6a468c77-3974-4ef0-a45c-a575c57c9b1b","470dab2c-b491-46e0-8b43-eb8ebc9f2d7b"]
    }
}
```

Now we can sent message in group

#### RPC Message send
```javascript
    socket.send(JSON.stringify({"method":40,"body":{"group_id":"95ba5e74-810d-49b4-943f-0b2753d688c9", "data":"Hello, how are you?"}}));
```
#### Received message sent event
```json
  {"type":21,"timestamp":1532786953431,"body":{"message_id":"f0b8dacc-74a6-4382-81a6-883e45d848a8"}}
```

#### On the user 2 side, received message event
```json
  {"type":20,"timestamp":1532786953431,"body":{"message_id":"f0b8dacc-74a6-4382-81a6-883e45d848a8","data":"Hello, how are you?"}}
``` 

User 2 should sent 2 rpc methods: message delivered and message read
#### RPC Message delivered
```js
  socket.send(JSON.stringify({"method":41,"body":{"message_id":"f0b8dacc-74a6-4382-81a6-883e45d848a8"}}));
```

User 1 received message delivered event
```json
  {"type":22,"timestamp":1532787352840,"body":{"message_id":"f0b8dacc-74a6-4382-81a6-883e45d848a8"}}
```

#### RPC Message read
```js
  socket.send(JSON.stringify({"method":42,"body":{"message_id":"f0b8dacc-74a6-4382-81a6-883e45d848a8"}}));
```

User 1 receive message read event
```json
{"type":23,"timestamp":1532787442736,"body":{"message_id":"f0b8dacc-74a6-4382-81a6-883e45d848a8"}}
```

User 2 typing start and typing end methods
```javascript
socket.send(JSON.stringify({"method":60,"body":{"group_id":"95ba5e74-810d-49b4-943f-0b2753d688c9"}}));
socket.send(JSON.stringify({"method":61,"body":{"group_id":"95ba5e74-810d-49b4-943f-0b2753d688c9"}}));
```

User 1 received typing start and typing end events
```json
{"type":40,"timestamp":1532787626297,"body":{"group_id":"95ba5e74-810d-49b4-943f-0b2753d688c9","user_id":"6a468c77-3974-4ef0-a45c-a575c57c9b1b"}}
```
```json
{"type":41,"timestamp":1532787626305,"body":{"group_id":"95ba5e74-810d-49b4-943f-0b2753d688c9","user_id":"6a468c77-3974-4ef0-a45c-a575c57c9b1b"}}
```

### Connect to server 
- connection to server for user 2
- sent rpc method for get all events on connection
```javascript
var socket = new WebSocket("ws://127.0.0.1:3000/b208dd54-3374-486b-9e27-9a19e4098f42");
socket.onopen = function() {
   socket.send(JSON.stringify({"method":20, "timestamp":new Date().getTime()}));
};
socket.onmessage = function(event) {
  console.log(event.data);
};
```

#### User 1 left group
```
curl -X PUT \
  http://127.0.0.1:3000/groups/95ba5e74-810d-49b4-943f-0b2753d688c9/left \
  -H 'authorization: bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzM5OTQ5MzksInVzZXJfaWQiOiI0NzBkYWIyYy1iNDkxLTQ2ZTAtOGI0My1lYjhlYmM5ZjJkN2IifQ.hE5NF5Dscdtn-Z2pdIyMBrXVeSBG69xeArr8DqoceWo' \
  -H 'content-type: application/json' \
  -d '{
	"user_id":"470dab2c-b491-46e0-8b43-eb8ebc9f2d7b"
}'
```

#### User 2 received event user group left
```json
  {"type":73,"timestamp":1532786538556,"body":{"group_id":"95ba5e74-810d-49b4-943f-0b2753d688c9","user_id":"470dab2c-b491-46e0-8b43-eb8ebc9f2d7b"}}
```

#### User 1 join group
```
curl -X PUT \
  http://127.0.0.1:3000/groups/95ba5e74-810d-49b4-943f-0b2753d688c9/join \
  -H 'authorization: bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzM5OTQ5MzksInVzZXJfaWQiOiI0NzBkYWIyYy1iNDkxLTQ2ZTAtOGI0My1lYjhlYmM5ZjJkN2IifQ.hE5NF5Dscdtn-Z2pdIyMBrXVeSBG69xeArr8DqoceWo' \
  -H 'content-type: application/json' \
  -d '{
	"user_id":"470dab2c-b491-46e0-8b43-eb8ebc9f2d7b"
}'
```
#### User 2 received event user group join
```json
  {"type":72,"timestamp":1532786538556,"body":{"group_id":"95ba5e74-810d-49b4-943f-0b2753d688c9","user_id":"470dab2c-b491-46e0-8b43-eb8ebc9f2d7b"}}
```