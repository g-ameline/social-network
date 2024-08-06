# instal alpine and golang image, get program data from repo, run container
FROM golang
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
# └ data/
#     ├ social-network.db
#     └ serve/
#         ├ all go stuff
#         └ database_server.go
# 
RUN mkdir ./social-network/
WORKDIR ./social-network/
# prepare back_end folder
# copy data folder  databse + server
RUN ls -lha
ADD ./data/social-network.db ./data/social-network.db
RUN ls -lha
ADD ./data/serve/ ./data/serve/
# move to server folder
WORKDIR ./data/serve
RUN ls -lha
EXPOSE 5000
# download dependencies
RUN go mod tidy
# RUN apk add --no-cache gcc musl-dev
# RUN CGO_ENABLED=1 GOOS=linux go build -ldflags=" -w -s" -o database_server
RUN ls -lha
# RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o data_server
# run the server
# CMD ["./database_server"]
CMD ["go","run","."]

