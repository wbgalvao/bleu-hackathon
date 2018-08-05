package bot

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	client "github.com/wbgalvao/bleu-hackathon/client"
	tb "gopkg.in/tucnak/telebot.v2"
)

const TOKEN_BOT = "676072443:AAGD_Ba7jDhlN3lIKij1Y4eZz7ImyktjQ_8"
const URL = "https://bleutrade.com/api/v2/"

var apiKey string
var apiSecret string
var cli client.Client

func Init() {
	b, err := tb.NewBot(tb.Settings{
		Token:  TOKEN_BOT,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}
	cli.BaseURL, _ = url.Parse(URL)
	cli.HttpClient = new(http.Client)

	b.Handle("/registerApiKey", func(m *tb.Message) {
		cli.APIKey = m.Payload
		b.Send(m.Sender, "Key registered")
	})

	b.Handle("/registerApiSecret", func(m *tb.Message) {
		cli.APISecret = m.Payload
		b.Send(m.Sender, "Api Secret registered")
	})

	b.Handle("/saldo", func(m *tb.Message) {
		b.Send(m.Sender, m.Payload)
		balances, err := cli.GetBalances()
		if err != nil {
			fmt.Println("Deu erro")
			b.Send(m.Sender, err)
		}
		for _, balance := range balances {
			fmt.Println(balance)

			b.Send(m.Sender, "Moeda: "+balance.Currency+"\nSaldo: "+balance.Balance)
		}
	})

	b.handle("/wallet", func(m *tb.Message) {

	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		// all the text messages that weren't
		// captured by existing handlers
		b.Send(m.Sender, m.Text)
	})

	fmt.Println("Bot started")
	b.Start()
}
