so you wanna run those three servers in containers
- bash ./run_docker.sh
hopefully it will be all good
when you are done enjoying our social network you can run
- stop_docker.sh
when you are done enjoying our social network **forever**
- clean_docker.sh
then if you have no other important image or container you wanna keep   

and if you really want to clean all docker stuff
- docker system prune -a

note : this is not state of art docker process, no multistage stuff no docker compose..
but as far as I can tell it works and stll use alpine and not ubuntu
cheers.