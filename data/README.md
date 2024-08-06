this is not databse helper package used by the server at runtime.  
social-network.sql :
- the actual databse used by the server at runtime
create :
- this is the script that create the database in create
- each script check if whatever it is supposed to add already exist
  - if not it will create it
  - if already there it pass
- this do what instruction call migration to my understanding you are welcome to argue with that
populate :
- populate the database with random values
  - it uses the same library than the server so it work as a test 
each of those are individual golang package and should be used as scripts either manually or by the server according to its own logic.

