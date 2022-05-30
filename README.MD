# WellNUS backend

## Constants

BACKEND_URL: localhost:8080

FRONTEND_URL: localhost:3000

## Features

### User

> #### Users Details
> 
>> Handles CRUD on users
> 
> #### User Routes
> 
>> ##### /user - GET
>> 
>> : Get all users
>> 
>> ##### /user - POST
>> 
>> : Create a new user and logins as that user
>> 
>> ##### /user/:id - GET
>> 
>> : Gets user with given id
>> 
>> ##### /user/:id - PATCH
>> 
>> : Update user with given id with updates, provided that user logged in
>> 
>> ##### /user/:id - DELETE
>> 
>> : Delete a user with given id, provided that user logged in

### Session

> #### Session Details
> 
>> Handles user authentication for registered users
> 
> #### Session Routes
> 
>> ##### /session - POST
>> 
>> : Login to appropriate user
>> 
>> ##### /session - DELETE
>> 
>> : Logout of any user

### Group

> #### Group Details
> 
>> Handles creation, and joining of groups.
>> 
>> Each group must have at least an owner present in said group.
> 
> #### Group Routes
> 
>> ##### /group - GET
>> 
>> : Get all groups
>> 
>> ##### /group - POST
>> 
>> : Creates a new group if user is loggedIn
>> 
>> ##### /group/:id - GET
>> 
>> : Get a group with given id
>> 
>> ##### /group/:id - POST
>> 
>> : Creates a join request in the wn_join_request table if user is loggedIn
>> 
>> ##### /group/:id - DELETE
>> 
>> : Leaves the group if user is not owner. Otherwise, deletes the group
>> 
>> ##### /group/:id - UPDATE
>> 
>> : Updates the information of the group if user is in given group

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

- user_id BIGINT REFERENCES wn_user(id)
- group_id BIGINT REFERENCES wn_group(id)
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
- [ ] CRUD on Groups
- [ ] Unit testing for Groups Querying
- [ ] Unit testing for Groups Handlers
- [ ] CRUD on Join_Request
- [ ] Implement Websockets for Groups