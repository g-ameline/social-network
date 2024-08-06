
here how the databse is designed
but first some rule
because of the way I wrote database functions all rows should have an id caled id even if other potential uniques identifiers are possibles


* main entities tables

users  
- id
- email
- password -> crypted
- first Name
- last Name
- birth
- avatar/Image (Optional)
- nickname (Optional)
- about Me (Optional)
- session_id
- last_activity
back session state will graved in database, server logic will decide if one session_id have been idle for too long and is expired not the client.


groups  
- id
- creator_id
- title
- description

notes  
- id
- sender
- predecessor
- text
note will be either post, chat message or group post|comment. its readership "label" will provide functonal disctinction.

* sub-groups

event of group

then the relations between each

followance user to user
- following state from a user toward another (one way)
  - requested
  - accepted
  - refused

readership (of notes)
- id
- type
  - exclusive (many to many link from note to chosen folowers) 
  - private (then all follower can see it target_id will be the followed author of the post)
  - public (then everyone can see it target_it could be the author or empty)
  - group (it is a message group the target_id will be the group_id)
- target_id
  - either group_id, user_id or empty see type
fetching cases :
- on profile component we need to fetch all posts from one user
  - we loop into readerships and keep the ones that match 
    - type = exclusise && target_id is current user_id ()
    - type = private and visited profile's user_id 
    - public  