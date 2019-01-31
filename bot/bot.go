package bot

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/wbgalvao/bleu-hackathon/balance"
	client "github.com/wbgalvao/bleu-hackathon/client"
	tb "gopkg.in/tucnak/telebot.v2"
)

const TOKEN_BOT = "676072443:AAGD_Ba7jDhlN3lIKij1Y4eZz7ImyktjQ_8"
const URL = "https://bleutrade.com/api/v2/"

var apiKey string
var apiSecret string

func NewClient(apiKey, apiSecret string) client.Client {
	var cli client.Client
	cli.BaseURL, _ = url.Parse(URL)
	cli.HttpClient = new(http.Client)

	cli.APIKey = apiKey
	cli.APISecret = apiSecret
	return cli
}

func Init() {
	b, err := tb.NewBot(tb.Settings{
		Token:  TOKEN_BOT,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}
	senderCache := make(map[int]client.Client)

	b.Handle("/registerApiKey", func(m *tb.Message) {
		var cli = senderCache[m.Sender.ID]
		cli.APIKey = m.Payload
		b.Send(m.Sender, "Key registered")
	})

	b.Handle("/registerApiSecret", func(m *tb.Message) {
		var cli = senderCache[m.Sender.ID]
		cli.APISecret = m.Payload
		b.Send(m.Sender, "Api Secret registered")
	})

	b.Handle("/setup", func(m *tb.Message) {
		splittedPayload := strings.Split(m.Payload, " ")
		apiKey := splittedPayload[0]
		apiSecret := splittedPayload[1]
		ncli := NewClient(apiKey, apiSecret)
		fmt.Println("new client mapped")
		fmt.Println(ncli)
		senderCache[m.Sender.ID] = ncli
		fmt.Println("Sender")
		fmt.Println(m.Sender.ID)

		b.Send(m.Sender, "Chave e segredo registrados!")
	})

	b.Handle("/saldo", func(m *tb.Message) {
		var balances []balance.Balance
		var err error
		cli := senderCache[m.Sender.ID]
		if m.Payload != "" {
			b.Send(m.Sender, "Seu saldo em "+m.Payload+"é:")
			balances, err = cli.GetBalances(m.Payload)
		} else {
			b.Send(m.Sender, "Você tem saldo nas seguintes moedas:")
			balances, err = cli.GetBalances()
		}
		if err != nil {
			fmt.Println("Deu erro")
			b.Send(m.Sender, err)
		}
		for _, balance := range balances {
			fmt.Println(balance)
			if n, err := strconv.ParseFloat(balance.Available, 32); n > 0 && err == nil {
				b.Send(m.Sender, "Moeda: "+balance.Currency+"\nSaldo: "+balance.Balance)
			}
		}
	})

	b.Handle("/wallet", func(m *tb.Message) {
		var cli = senderCache[m.Sender.ID]
		fmt.Println("sender wallet")
		fmt.Println(m.Sender)
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
		var cli = senderCache[m.Sender.ID]
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
		var cli = senderCache[m.Sender.ID]
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

	b.Handle("/saque", func(m *tb.Message) {
		var cli = senderCache[m.Sender.ID]
		b.Send(m.Sender, m.Payload)
		splittedPayload := strings.Split(m.Payload, " ")
		currency := splittedPayload[0]
		quantity := splittedPayload[1]
		walletDest := splittedPayload[2]
		result, err := cli.Withdraw(currency, quantity, walletDest)
		if err != nil {
			fmt.Println("Deu erro")
			b.Send(m.Sender, fmt.Sprintf("%v", err))
		}
		if result {
			b.Send(m.Sender, fmt.Sprintf("Saque efetuado!"))
		} else {
			b.Send(m.Sender, fmt.Sprintf("Problema no saque"))
		}
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		b.Send(m.Sender, m.Text)
	})

	fmt.Println("Bot started")
	b.Start()
}
