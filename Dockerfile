FROM gliderlabs/alpine:3.2
ADD nativerw config.json /
EXPOSE 8080

CMD /nativerw -mongos=$MONGO_ADDRESSES config.json

