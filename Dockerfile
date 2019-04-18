FROM golang:latest

# Set apps working directory
WORKDIR /app

# Set an env var that matches $GOPATH
ENV SRC_DIR=/go/src/github.com/ashkarin/ashkarin-api-test

# Get dependencies
RUN go get github.com/gorilla/mux
RUN go get gopkg.in/mgo.v2
RUN go get github.com/sirupsen/logrus
# Tests
RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega


# Copy the local files to the container's workspace
ADD /pkg $SRC_DIR/pkg
ADD /internal $SRC_DIR/internal
ADD /configs $SRC_DIR/configs
ADD /cmd $SRC_DIR/cmd

# Build and copy to compiled directory
RUN cd $SRC_DIR/cmd/server; go build; cp server /app/

RUN echo '#!/bin/bash \n\
if [[ "$TEST" == "true" ]] \n\
then \n\
    echo "TEST: $TEST"\n\
    cd $SRC_DIR\n\
    ginkgo -r -v \n\
else \n\
    echo "SERVER"\n\
    /app/server\n\
fi' > /app/entry.sh && chmod +x /app/entry.sh

RUN cat /app/entry.sh

# Set entrypoint
#ENTRYPOINT ./server -config $SRC_DIR/configs/config.json
ENTRYPOINT ./entry.sh

EXPOSE 80 8080