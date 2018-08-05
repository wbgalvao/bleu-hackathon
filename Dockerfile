FROM golang


COPY . $GOPATH/src/github.com/wbgalvao/bleu-hackathon
WORKDIR $GOPATH/src/github.com/wbgalvao/bleu-hackathon

RUN go get "gopkg.in/tucnak/telebot.v2"

RUN go build .

CMD ["./bleu-hackathon"]

