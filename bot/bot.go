package bot

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
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

	b.Handle("/setup", func(m *tb.Message) {
		splittedPayload := strings.Split(m.Payload, " ")
		apiKey := splittedPayload[0]
		apiSecret := splittedPayload[1]
		cli.APIKey = apiKey
		cli.APISecret = apiSecret
		b.Send(m.Sender, "Chave e segredo registrados!")
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

	b.Handle("/wallet", func(m *tb.Message) {
		b.Send(m.Sender, m.Payload)
		result, err := cli.GetBalances("BTC")
		if err != nil {
			fmt.Println("Deu erro")
			b.Send(m.Sender, err)
		}
		for _, balance := range result {
			if balance.Currency == "BTC" {
				b.Send(m.Sender, fmt.Sprintf("Aqui está o endereço da sua wallet: %s", balance.CryptoAddress))
			}
		}
	})

	b.Handle("/buylimit", func(m *tb.Message) {
		b.Send(m.Sender, m.Payload)
		splittedPayload := strings.Split(m.Payload, " ")
		market := splittedPayload[0]
		quantity := splittedPayload[1]
		result := make(map[string]string)
		result, err := cli.BuyLimit(market, quantity)
		if err != nil {
			fmt.Println("Deu erro")
			b.Send(m.Sender, err)
		}
		b.Send(m.Sender, fmt.Sprintf("Compra efetuada! Identificação da transação: %s", result["orderid"]))
	})

	b.Handle("/selllimit", func(m *tb.Message) {
		b.Send(m.Sender, m.Payload)
		splittedPayload := strings.Split(m.Payload, " ")
		market := splittedPayload[0]
		quantity := splittedPayload[1]
		result := make(map[string]string)
		result, err := cli.SellLimit(market, quantity)
		if err != nil {
			fmt.Println("Deu erro")
			b.Send(m.Sender, err)
		}
		b.Send(m.Sender, fmt.Sprintf("Venda efetuada! Identificação da transação: %s", result["orderid"]))
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		// all the text messages that weren't
		// captured by existing handlers
		b.Send(m.Sender, m.Text)
	})

	fmt.Println("Bot started")
	b.Start()
}
