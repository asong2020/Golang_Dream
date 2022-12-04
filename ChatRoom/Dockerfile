FROM library/golang:1.14

ENV GOROOT=/usr/local/go
ENV GOPATH=/root/go/
ENV PATH=$GOPATH/bin/:$PATH
ENV APP_DIR $GOPATH/src/asong.cloud/ChatRoom
RUN go get github.com/astaxie/beego && go get github.com/beego/bee & go get github.com/go-sql-driver/mysql
WORKDIR $APP_DIR
ADD . $APP_DIR
EXPOSE 8080

CMD ["bee", "run"]