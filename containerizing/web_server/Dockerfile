FROM node:alpine
# FROM node:alpine
# recreate necesary repo structure from
# ┘ ┐ ┌ └ ┼ ─ ├ ┤ ┴ ┬  │
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
# 
# to reproduce :
# social-network (dockerfile start from here)
# └ presentation/
#     └ Caddyfile
RUN mkdir ./social-network/
WORKDIR ./social-network/
ADD ./presentation/ ./presentation/
WORKDIR ./presentation/
# prepare back_end presentation folder
# copy front-end files
RUN ls -lha 
RUN apk add caddy
RUN ls -lha
RUN caddy adapt
EXPOSE 3000
CMD ["caddy","run"]

