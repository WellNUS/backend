# WellNUS backend

## Details
BACKEND_URL: localhost:8080\
FRONTEND_URL: localhost:3000\

## Routes
### Users
/user - GET\
: Get all users\
/user - POST\
: Create a new user and logins as that user\
/user/:id - GET\
: Gets user with given id\
/user/:id - PATCH\
: Update user with given id with updates, provided that user logged in\
/user/:id - DELETE\
: Delete a user with given id, provided that user logged in\

### Sessions
/session - POST\
: Login to appropriate user\
/session - DELETE\
: Logout

### Groups
/group - GET\
: Get all groups
/group/:id - GET\
: Get group with given id\

## Things to do
- [x] CRUD on Users
- [x] CRUD on Sessions
- [x] Unit testing for Users Querying
- [x] Unit testing for Users Handlers
- [x] Unit testing for Sessions
- [ ] CRUD on Groups
- [ ] Implement Websockets for Groups
- [ ] Unit testing for Groups Querying
- [ ] Unit testing for Groups Handlers