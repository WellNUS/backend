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

### User

> #### Users Details
> 
>> Handles CRUD on users
>> 
>> User = { id, first_name, last_name, gender, faculty, email, user_role, password, password_hash }
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
>>> Description : Gets user with given id
>>>
>>> Request Body : None
>>>
>>> Response Body: User
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
>>> Description : Get a group with given id
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
>> JoinRequest = { id, user_id, group_id, request_status }
>
> #### Join Routes
>
>> ##### /join - GET
>>
>>> Description : Get all join requests sent by user or directed at user. Query params used to indicate which request=SENT)
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
>>> Description : Gets join request with the given id
>>>
>>> Request Body : None
>>>
>>> Response Body : JoinRequest
>>
>> ##### /join/:id - PATCH
>>
>>> Description : Approve or reject a particular join request if the user is the owner of the particular group
>>>
>>> Request Body : { approve }
>>>
>>> Response Body : JoinRequest
>>
>> #### /join/:id - DELETE
>>
>>> Description : Deletes the join request if the user is the owner of the request.
>>>
>>> Request Body : None
>>>
>>> Response Body : { id }

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

## Tables

Table creation with POSTGRESQL

### wn_user

- id BIGSERIAL NOT NULL PRIMARY KEY
- first_name VARCHAR(32) NOT NULL
- last_name VARCHAR(32) NOT NULL
- gender VARCHAR(1) NOT NULL
- faculty VARCHAR(10) NOT NULL
- email VARCHAR(128) NOT NULL
- user_role VARCHAR(10) NOT NULL
- password_hash VARCHAR(100) NOT NULL
- unique(email),
- check(first_name != ''),
- check(last_name != ''),
- check(password_hash != ''),
- check(gender IN ('M', 'F')),
- check(faculty IN ('CHS', 'BUSINESS', 'COMPUTING', 'DENTISTRY', 'CDE', 'LAW', 'MEDICINE', 'NURSING', 'PHARMACY', 'MUSIC')),
- check(email LIKE '%@u.nus.edu'),
- check(user_role IN ('MEMBER', 'VOLUNTEER', 'COUNSELLOR'))

### wn_group

- id BIGSERIAL NOT NULL PRIMARY KEY
- group_name VARCHAR(64) NOT NULL
- group_description VARCHAR(256) NOT NULL
- category VARCHAR(7) NOT NULL
- owner_id BIGINT REFERENCES wn_user(id) NOT NULL
- check(group_name != ''),
- check(category IN ('COUNSEL', 'SUPPORT'))

### wn_user_group

- user_id BIGINT REFERENCES wn_user(id) ON DELETE CASCADE
- group_id BIGINT REFERENCES wn_group(id) ON DELETE CASCADE
- unique(user_id, group_id)

### wn_join_request

- id BIGSERIAL NOT NULL PRIMARY KEY,
- user_id BIGINT REFERENCES wn_user(id),
- group_id BIGINT REFERENCES wn_group(id),
- request_status VARCHAR(8) NOT NULL,
- unique(user_id, group_id)
- check(request_status IN ('PENDING', 'REJECTED', 'APPROVED'))

### wn_match_request
- id BIGSERIAL NOT NULL PRIMARY KEY,
- user_id BIGINT REFERENCES wn_user(id),
- cross_gender BOOLEAN NOT NULL,
- cross_faculty BOOLEAN NOT NULL,
- size SMALLINT NOT NULL, 
- mbti VARCHAR(4) NOT NULL,
- request_status VARCHAR(8),
- unique(user_id),
- check(size <= 8),
- check(mbti IN ('INFJ', 'INTP', 'ENTJ', 'ENTP', 'INFJ', 'INFP', 'ENFJ', 'ENFP', 'ISTJ', 'ISFJ', 'ESTJ', 'ESFJ', 'ISTP', 'ISFP', 'ESTP', 'ESFP'))
- check(request_status IN ('PENDING', 'REJECTED', 'APPROVED'))

### wn_counsel_request
- id BIGSERIAL NOT NULL PRIMARY KEY,
- user_id BIGINT REFERENCES wn_user(id),
- counsel_title VARCHAR(256) NOT NULL,
- counsel_description VARCHAR(256) NOT NULL,
- request_status VARCHAR(8) NOT NULL,
- unique(user_id),
- check(counsel_title != ''),
- check(counsel_description != ''),
- check(request_status IN ('PENDING', 'REJECTED', 'APPROVED'))

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
- [ ] Rewrite README.MD
- [ ] Implement Websockets for Groups