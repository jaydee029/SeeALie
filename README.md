# SEEALIE

A cli based Chat application, the chat server is developed demonstrating the microservice architecture demonstrating various backend technologies, it boasts of three services a chat service , an auth/user service, an email service.
Apart from that a cli client which would allow you to use the application. (Project Under development)

### Schema

#### User service
- users- id , username, email, password, created_at
- sessions- session_id, user_id, jwt_token, expires_at fk userid to userid in users
- revoked- token, revoked_at

#### Chat service
- friends- followed, followed_by, room_id, connected_at
- connections- connection_id, requested_by, recieved_by, created_at, approved_at, status
  
