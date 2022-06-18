# WellNUS backend

## Installation and Setup

1. Get docker desktop on your machine
2. Get [migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) dependency
2. clone repository to you machine and cd into the folder
3. To setup the database and tables, run the following

```
    $ make postgres
    $ make createdb
    $ make migrateup
```

Note: docker postgres will be using port 5432 on localhost. 

Ensure that given port is available before running `make postgres`. 

Otherwise, you can change the port in 'Makefile'.

4. To start run `go run .`

## Constants

BACKEND_URL: localhost:8080

FRONTEND_URL: localhost:3000

## Features

### Notes to connect frontend to backend

All Request made should have **credentials: "include"** so that cookies can be accessed and modified by backend

### User

> #### Users Details
> 
>> Handles CRUD on users
>> 
>> User = { id, first_name, last_name, gender, faculty, email, user_role, password, password_hash }
>>
>> UserWithGroups = { user: User, groups: Group[] }
> 
> #### User Routes
> 
>> ##### /user - GET
>>
>>> Description: Get all users
>>> 
>>> Request Body: None
>>> 
>>> Response Body: User[]
>>
>> ##### /user - POST
>> 
>>> Description : Create a new user and logins as that
>>>
>>> Request Body: { first_name, last_name, gender, faculty, email, user_role, password }
>>>
>>> Response Body: User
>>
>> ##### /user/:id - GET
>> 
>>> Description : Gets user with given id and the user's corresponding groups 
>>>
>>> Request Body : None
>>>
>>> Response Body: UserWithGroups
>>
>> ##### /user/:id - PATCH
>> 
>>> Description : Update user with given id with updates, if the user is logged in
>>>
>>> Request Body : { first_name?, last_name?, gender?, faculty?, email?, user_role?, password? }
>>>
>>> Response Body : User
>>
>> ##### /user/:id - DELETE
>> 
>>> Description : Delete a user with given id, if the user is logged in
>>>
>>> Request Body : None
>>>
>>> Response Body : { id }

### Session

> #### Session Details
> 
>> Handles user authentication for registered users
> 
> #### Session Routes
> 
>> ##### /session - POST
>> 
>>> Description : Login to appropriate user
>>>
>>> Request Body : { email, password }
>>> 
>>> Response Body : { logged_in, user: User }
>> 
>> ##### /session - DELETE
>> 
>>> Description : Logout of any user
>>> 
>>> Request Body : None
>>> 
>>> Response Body : { logged_in, user: User }

### Group

> #### Group Details
> 
>> Handles creation, and joining of groups.
>> 
>> Each group must have at least an owner present in said group.
>>
>> Group = { id, group_name, group_description, category, owner_id } 
>>
>> GroupWithUsers = { group: Group, users: User[] }
> 
> #### Group Routes
> 
>> ##### /group - GET
>> 
>>> Description : Get all groups that the user is in if the user is logged in.
>>> 
>>> Request Body : None
>>>
>>> Response Body : Group[]
>> 
>> ##### /group - POST
>> 
>>> Description : Creates a new group if user is logged in
>>>
>>> Request Body : { group_name, group_description?, category }
>>>
>>> Response Body : GroupWithUsers
>> 
>> ##### /group - DELETE
>>
>>> Description : Leave all groups, transferring ownership of group to a random user in each group if necessary. If it is the last user of a group, the group is deleted
>>> 
>>> Request Body : None
>>>
>>> Response Body : GroupWithUsers[]
>>
>> ##### /group/:id - GET
>> 
>>> Description : Get a group with given id and the users within the group
>>>
>>> Request Body : None
>>>
>>> Response Body : GroupWithUsers
>> 
>> ##### /group/:id - DELETE
>> 
>>> Description : Leaves the group, transferring ownership of group to a random user in group if necessary. If the last user leaves, the group is deleted.
>>>
>>> Request Body : None
>>>
>>> Response Body : GroupWithUsers
>>
>> ##### /group/:id - UPDATE
>> 
>>> Description : Updates the information of the group if user is the owner of the given group
>>> 
>>> Request Body : { group_name?, group_description?, category?, owner_id? } 
>>>
>>> Response Body : Group

### Join

> #### Join Details
>
>> Handles CRUD on group join requests
>>
>> JoinRequest = { id, user_id, group_id }
>>
>> LoadedJoinRequest = { join_request: JoinRequest, user: User, group: Group }
>>
>> JoinRequestRespond = { approve }
>
> #### Join Routes
>
>> ##### /join - GET
>>
>>> Description : Get all join requests sent by user or directed at user. Query params used to indicate what join request to retrieve.
>>>
>>> Query Params:
>>> - ?request=SENT     will get all join request sent by user
>>> - ?request=RECEIVED will get all join request directed at user
>>> - no query params   will get all join request sent by and directed at user
>>> 
>>> Request Body : None
>>> 
>>> Response Body : JoinRequest[]
>>
>> ##### /join - POST
>>
>>> Description : Creates new join request to a certain group if the user is logged in.
>>>
>>> Request Body : { group_id }
>>> 
>>> Response Body : JoinRequest
>>
>> ##### /join/:id - GET
>>
>>> Description : Gets loaded join request with the given id
>>>
>>> Request Body : None
>>>
>>> Response Body : LoadedJoinRequest
>>
>> ##### /join/:id - PATCH
>>
>>> Description : Approve or reject a particular join request if the user is the owner of the particular group
>>>
>>> Request Body : JoinRequestRespond
>>>
>>> Response Body : JoinRequestRespond
>>
>> ##### /join/:id - DELETE
>>
>>> Description : Deletes the join request if the user is the owner of the request.
>>>
>>> Request Body : None
>>>
>>> Response Body : { id }

### Message

> #### Message Details
>
>> Handles the retrieval of past messages in groups. 
>>
>> The API takes in query params to indicate the limit and the latest message to load in. This is to enable dynamic loading of messages so that not all the chat history is loaded in at once. Refer to /template/chat/chat.html for reference on how to do this
>>
>> Message = { user_id, group_id, time_added, msg }
>>
>> MessagePayload = { tag=0, sender_name, group_name, message: Message }
>>
>> MessagesChunk = { earliest_time, latest_time, message_payloads: message_payload[] }
>
> #### Message Routes
>
>> ##### /message/:id - GET
>>
>>> Description : Gets messages chunk sent in the group of given id in ascending time sent.
>>>
>>> Query Params:
>>> - ?limit=(int) : Specify how many messages to load in. If unspecified or <= 0, will load all past messages
>>> - ?latest=(time) : Specify the latest time sent of the retrieved messages. If unspecified, will take current time as latest time. Note: Time specified must be in RFC3339Nano format. Example: "2006-01-02T15:04:05.999999999Z07:00"
>>>
>>> Request Body : None
>>>
>>> Response Body : MessagesChunk

### Websocket for chats

> #### Websocket details
>
>> Allows client to obtain a websocket connection to the server.
>>
>> The following is an example from /template/chat/chat.html on how to connect to the websocket with vanilla javascript. Head over to /template/chat/chat.html to see the full code
```
conn = new WebSocket("ws://" + document.location.host + "/ws/" + groupID);
conn.onclose = function (evt) {
    var item = document.createElement("div");
    item.innerHTML = "<b>Connection closed.</b>";
    appendLog(item);
};
conn.onmessage = function (evt) {
    var payload = JSON.parse(evt.data)
    handlePayload(payload)
};
```
>> Sending data to server:
>>
>>> The data sent to the server through the **conn** instance has a format of **Message** = { user_id, group_id, time_added and msg }.
>>
>> Receiving data from server:
>>
>>> The data received from the server through the **conn** instance is what we will call a **Payload**.
>>>
>>> A Payload always comprise of a **tag** field to indicate what kind of data has been passed through.
>>>
>>> Tag Reference sheet:
>>> - tag == 0 **MessagePayload**
>>>     - **MessagePayload** = { tag, sender_name, group_name, message: Message }
>>>         - If the **message.user_id == -1**, the message a Server Message and is not saved on database.
>>> - tag == 1 **ChatStatusPayload**
>>>     - **ChatStatusPayload** = { tag, group_id, group_name, sorted_in_chat_members, sorted_online_members, sorted_offline_members }
>>>         - Members of a group are in only 1 of 3 states
>>>             - in chat : Member is connected through websocket and is on chat page of the given group
>>>             - online : Member is connected through websocket and is on chat page of some other group
>>>             - offline : Member is not connected through websocket
>>>         - Union of **sorted_in_chat_members**, **sorted_online_members**, and **sorted_offline_members** form all group members of given group
>>> These data sent from the server is sufficient to create the following chat features:
>>> - Live messages from users
>>> - Live messages from server
>>> - Notifications of users joining the chat
>>> - Notifications of users leaving the chat
>>> - Live updates of what users are online in a given group
>>> - Live updates of what users are in chat of the given group
>> 
>> Refer to /templates/chat/chat.html for reference on how to implement the chat feature
>
> #### Websocket Routes
>
>> ##### /ws/:group_id - GET
>>
>>> Description : Create a websocket connection with the url of "ws://(domain_name)/ws/(group_id)" which will return a connection object. Incoming payloads and outgoing messages will come and go through this connection object.
>>> 
>>> Request Body : None
>>>
>>> Response Body : None

### Routes past this point have not been implemented
### Match

> #### Match Details
> 
>> All matches are made when a certain number of match requests are present.
>> 
>> Current matching algorithm:
>> 
>> 1. Randomly choose a match request
>> 2. Set requestee as the owner of the group
>> 3. Greedily satisfy a compatibility comparative function with other requests
>> 4. Form group users and remove their requests
>> 5. Repeat steps 1 to 5 till either all requests are fulfiled or there are not sufficient users
>> 
>> Compatability function:
>> 
>> - Compatibility function will give a score from 0 - 100 between 2 match request.
>> - Function will consider gender, faculty, preferred group size, and MBTI type
> 
> #### Match Routes
> 
>> ##### /match - GET
>> 
>> : Gets all match request
>> 
>> ##### /match - POST
>> 
>> : Creates a new match request if user is loggedIn
>> 
>> ##### /match/:id - GET
>> 
>> : Get match request with given id
>> 
>> ##### /match/:id - DELETE
>> 
>> : Delete match request with given id if user is the owner of match request
>> 
>> ##### /match/:id - UPDATE
>> 
>> : Update match request with given id if user is the owner of match request

### Counsel

> #### Counsel Details
> 
>> Similar to a job listing site. Users whose are user_role = COUNSEL will be able to select to take on various request.
>> 
>> Each request a time limit of 3 days. Request approaching time limit will be highlighted for volunteer counsellors to take on.
>> 
>> After request time limit, a non-volunteer paid counsellor will be assigned to take on request at their discretion.
> 
> #### Counsel Routes
> 
>> ##### /counsel - GET
>> 
>> : Gets all counsel request
>> 
>> ##### /counsel - POST
>> 
>> : Creates a new counsel request if user is loggedIn
>> 
>> ##### /counsel/:id - GET
>> 
>> : Gets counsel request with given id
>> 
>> ##### /counsel/:id - POST
>> 
>> : Assigns counsel request to user if user is a counsellor or volunteer
>> 
>> ##### /counsel/:id - DELETE
>> 
>> : Deletes a counsel request if user is the owner of the request
>> 
>> ##### /counsel/:id - UPDATE
>> 
>> : Updates a counsel request if user is the owner of the request

## Want to test this backend API on a GUI?

Test some of the features by accessing "http://(domain_name)/testing" after running the server. 

Note that not every API implemented is implemented in the above testing GUI interface (coming soon if there is time)

## Things to do
- [x] CRUD on Users
- [x] CRUD on Sessions
- [x] Unit testing for Users Querying
- [x] Unit testing for Users Handlers
- [x] Unit testing for Sessions
- [x] CRUD on Groups
- [x] Unit testing for Groups Handlers
- [x] CRUD on Join_Request
- [x] Unit testing for Join_Request
- [x] Configure database migration
- [x] Restructured file structure
- [ ] Support forget password feature
- [x] Rework join request for better usability (prevent request to alrdy joined group, allow automatic deletion)
- [x] Rewrite README.MD
- [x] Implement Websockets for Groups
- [x] Add dynamic loading for past messages
- [x] Added notifications to testing chat
- [x] Complete restructuring of code for hopefully better scalability
- [x] Generalised the websocket connection to enable chat status data and messages as payloads
- [x] Secured message loading with cookie
- [x] Added session authorisation with session database
- [x] Rewrite testcases for change