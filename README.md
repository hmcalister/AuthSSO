# AuthSSO

WARNING: This project is still a work in progress, and is NOT currently ready for deployment.

This project is an implementation of a single-sign-on system, where a single HTTP server will handle registering users, validating login attempts, and giving auth tokens / sessions. The main focus of this project is an investigation into best practices for storing confidential information such as passwords in a secure format, as well as creating a centralized system to handle user sessions in a web-application context.

## TODO

- API versioning
- Secret key rotation
- Refresh JWT

## Technologies Used

### Authentication 
argon2ID

### Session Management
JWTAuth

### Data Storage
SQLite
SQLC

### Networking
Go stdlib
Chi Router

### Logging 
Zerolog
Lumberjack

### Other
GNU Make
