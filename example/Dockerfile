FROM registry.matchvs.com/language/alpine:1.0.0
WORKDIR /gameServer-go
COPY . /gameServer-go
RUN ["chmod", "+x", "gameserver_go"]
ENTRYPOINT ["./gameserver_go"]