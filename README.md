
Hi, 
this is my social-network project.
That project have been done as part of mandatory task of the kood/johvi curriculum
key features are :
- following system
- group creation
- group's event creation 
- notifications
...

esthetism was not a priority , I hope it will be functional enough for you though.
constructive feed back is welcome indeed.

## how to run it either ##
- at the root : 
~run `bash start_servers.sh` (assuming you have sqlite and ***caddy*** installed~ that one run program in background, it can be confusing
-run 
 - `bash start_database_server.sh` assuming you have golang and sqlite installed
 - `bash start_logic_server.sh` 
 - `bash start_presentation_server.sh` (assuming you have ***caddy*** installed
-same thing but starting easch 3 servers manually in 
 - data/ `go run .`
 - logic/ `go run .`
 - presentation/ `caddy run`
- in containerizing/ run `bash run_dockers.sh`
- in containerizing/
 - in data_server/ run `bash run_docker.sh`
 - in app_server/ run `bash run_docker.sh`
 - in web_server/ run `bash run_docker.sh`
(in containerizing/ `bash clean_dockers.sh` to clean up after)
* please note that the database server might take a long time to be running (could not fix that), from docker or not, sorry for inconvenience *
### then connect to http://localhost:3000/ with your browser ###

## design details ## 
- chat is not persistent, it is stored locally on server
 - this does not break any audit spec AFAIK
- project is split in 3 layers
 - presentation (just a caddy reverse proxy server that serve SPA and echo the rest to application server)
  - port 3000
 - logic (where app server lives, do all the logic reponding)
  - port 8000
 - data (where database is stored and repond to query from app server)
  - port 5000
- database migration stuff 
 - not sure what they want us to od but :
  - there is a create/ module that is a golang script that create db from scratch
  - there is a populate/ module that is a golang script that add dummy entities to the created db
  - I assume that whatever we are supposed to do, the above script have at least the backbone to do it
- data is cached in the app server, meaning that :
 - data might be fetched from database only one time,
 - only statement are systematically send to database server
 - data on app layer are stored in non concurent safe map
  - high request rate will involve issue untill channel or mutex is added
- **HTMX** framwork is used for front end rendering through websocket
 - this is server side rendering (SSR)
 - all is handled in the logic layer, there is no logic in the presentation layer
 - a little homemade library is used to create html/htmx fragment programmatically

## (simplified) repo structure ##
 ./social-network  
├──  data/  
│  ├── social-network.db/  
│  ├──  create/  
│  ├──  populate/  
│  └──  serve/  
│       ├── all golang modules stuff ...  
│       └── database_server  
├──  logic/
│  ├── all golang modules stuff ...  
│  └── application_server.go  
├──  presentation/  
│  └── Caddyfile  *reverse proxy* ~ web server
├──  README.md  
├── containerizing/  
│     ├── data_server/  
│     │    ├── Dockerfile  
│     │    └── run_docker.sh
│     ├── app_server/  
│     │    ├── Dockerfile  
│     │    └── run_docker.sh
│     ├── web_server/  
│     │    ├── Dockerfile  
│     │    └── run_docker.sh
│     ├── clean_dockers.sh
│     └── run_dockers.sh
└── run_servers  
    
