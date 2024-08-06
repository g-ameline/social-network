# instal alpine and golang image, get program data from repo, run container
FROM golang:alpine
# recreate necesary repo structure from
# social-network (dockerfile start from here)
# ├ presentation/
# │    └ Caddyfile
# ├ logic/
# │    ├ all go stuff
# │    └ application_server.go
# └ data/
#     ├ social-network.db
#     └ serve/
#         ├ all go stuff
#         └ database_server.go

# to reproduce :
# social-network (dockerfile start from here)
# └ logic/
#    ├ all go stuff
#    └ application_server.go
RUN mkdir ./social-network/
WORKDIR ./social-network/
ADD ./logic/ ./logic/
WORKDIR ./logic/
# download dependencies
EXPOSE 8000
RUN ls -lha
RUN go build -o app_server
CMD ["./app_server"]

