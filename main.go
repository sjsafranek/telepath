package main

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"log"
	"errors"

	"github.com/pelletier/go-toml"
	"gopkg.in/gomail.v2"
)

/*

usage

go run main.go -config gmail.toml -attachment myfile.csv -subject 'here is myfile' user@email.com

*/

const (
	DEFAULT_CONFIG_FILE string = "config.toml"
	DEFAULT_ATTACHMENT  string = ""
	DEFAULT_SUBJECT     string = ""
	DEFAULT_MESSAGE     string = ""
)

var (
	CONFIG_FILE = DEFAULT_CONFIG_FILE
	ATTACHMENT  = DEFAULT_ATTACHMENT
	SUBJECT     = DEFAULT_SUBJECT
	MESSAGE     = DEFAULT_MESSAGE
)

type Config struct {
	Title    string `toml:"title"`
	SmtpHost string `toml:"smpt_host"`
	SmtpPort int    `toml:"smpt_port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Name     string `toml:"name"`
	Email    string `toml:"email"`
}

func main() {
	flag.StringVar(&CONFIG_FILE, "config", DEFAULT_CONFIG_FILE, "config file")
	flag.StringVar(&ATTACHMENT, "attachment", DEFAULT_ATTACHMENT, "attachement")
	flag.StringVar(&SUBJECT, "subject", DEFAULT_SUBJECT, "subject")
	flag.StringVar(&MESSAGE, "message", DEFAULT_MESSAGE, "message")
	flag.Parse()

	if "" == SUBJECT && "" == MESSAGE {
		log.Fatal(errors.New("subject and message cannot be blank"))
	}

	log.Println("opening config file")
	var conf Config
	b, err := ioutil.ReadFile(CONFIG_FILE)
	if nil != err {
		log.Fatal(err)
	}
	err = toml.Unmarshal(b, &conf)
	if nil != err {
		log.Fatal(err)
	}

	log.Println("collecting recievers")
	recievers := flag.Args()

	// build email message
	log.Println("building email message")
	m := gomail.NewMessage()
	m.SetHeader("From", conf.Email)
	m.SetHeader("To", recievers...)
	m.SetAddressHeader("Cc", conf.Email, conf.Name)
	m.SetHeader("Subject", SUBJECT)
	m.SetBody("text/html", MESSAGE)
	if "" != ATTACHMENT {
		m.Attach(ATTACHMENT)
	}

	// connect to smtp
	log.Println("connecting to smtp")
	d := gomail.NewDialer(conf.SmtpHost, conf.SmtpPort, conf.Username, conf.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email
	log.Println("sending email")
	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
	}

	log.Println("email has been sent")
}
