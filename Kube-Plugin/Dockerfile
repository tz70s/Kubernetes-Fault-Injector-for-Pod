FROM golang

RUN go get github.com/tools/godep && \
	go get k8s.io/client-go/... && \
	cd $GOPATH/src/k8s.io/client-go && \
	git checkout -b v2.0.0 && \
	godep restore ./...

RUN mkdir /app
ADD . /app/
WORKDIR /app


RUN go build -o kube-injector-plugin .
CMD ["/app/kube-injector-plugin"]
