package courier

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	retCode := m.Run()
	os.Exit(retCode)
}

func TestSendsEmail(t *testing.T) {
	cfg := &sp.Config{
		BaseUrl:    "https://api.sparkpost.com",
		ApiKey:     os.Getenv("SPARKPOST_KEY"),
		ApiVersion: 1,
	}
	var client sp.Client
	err := client.Init(cfg)
	if err != nil {
		log.Fatalf("SparkPost client init failed: %s\n", err)
	}

	ioutil.WriteFile("/tmp/gocourier-test.txt", []byte("Attachment Test"), 0666)

	f, err := os.Open("/tmp/gocourier-test.txt")

	if err != nil {
		log.Fatalf("Trouble creating test file: %s\n", err)
	}

	c := SparkPostCourier{
		client: client,
	}

	e := Email{
		To: []Address{
			Address{
				Email: os.Getenv("COURIER_RECIPIENT"),
				Name:  "Receiver",
			},
		},
		From: Address{
			Email: os.Getenv("SPARKPOST_SENDER"),
			Name:  "Sender",
		},
		Subject: "Test From Go Land",
		Headers: map[string]string{
			"X-Test-Header": "test",
		},
		Attachments: []Attachment{
			FileAttachment{
				File: f,
			},
		},
		Content: SimpleContent{
			Text: "Simple Text",
			HTML: "<strong>Simple HTML</strong>",
		},
	}

	id, err := c.Send(e)

	if err != nil {
		t.Error(err)
	} else {
		t.Log("Transmission ID is: " + id)
	}
}

func TestSendsTemplatedEmail(t *testing.T) {
	// apiKey := os.Getenv("SPARKPOST_API_KEY")
	cfg := &sp.Config{
		BaseUrl:    "https://api.sparkpost.com",
		ApiKey:     os.Getenv("SPARKPOST_KEY"),
		ApiVersion: 1,
	}
	var client sp.Client
	err := client.Init(cfg)
	if err != nil {
		log.Fatalf("SparkPost client init failed: %s\n", err)
	}

	c := SparkPostCourier{
		client: client,
	}

	ioutil.WriteFile("/tmp/gocourier-test.txt", []byte("Attachment Test"), 0666)

	f, err := os.Open("/tmp/gocourier-test.txt")

	if err != nil {
		log.Fatalf("Trouble creating test file: %s\n", err)
	}

	e := Email{
		To: []Address{
			Address{
				Email: os.Getenv("COURIER_RECIPIENT"),
				Name:  "Receiver",
			},
		},
		From: Address{
			Email: os.Getenv("SPARKPOST_SENDER"),
			Name:  "Sender",
		},
		Subject: "Templated Test From Go Land",
		Headers: map[string]string{
			"X-Test-Header": "test",
		},
		Content: TemplatedContent{
			TemplateID: os.Getenv("SPARKPOST_TEMPLATE_ID"),
			SubstitutionData: map[string]string{
				"html": `<strong>Templated HTML</strong>`,
				"text": "Templated Text",
			},
		},
		Attachments: []Attachment{
			FileAttachment{
				File: f,
			},
		},
	}

	id, err := c.Send(e)

	if err != nil {
		t.Error(err)
	} else {
		t.Log("Transmission ID is: " + id)
	}
}
