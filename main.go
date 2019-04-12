package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
)

/*

usage

go run main.go -config gmail.toml -attachment myfile.csv -subject 'here is myfile' user@email.com

*/

const (
	DEFAULT_CONFIG_FILE string = "telepath.toml"
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

func getConfig(filename string) (Config, error) {
	log.Println("opening config file")
	var conf Config
	err := conf.Fetch(CONFIG_FILE)
	if nil != err {

		ir := InputReader{}
		fmt.Println("Setup email config")
		fmt.Println("------------------")
		conf.Title = "telepath"
		conf.Username, _ = ir.Read("username: ")
		conf.Password, _ = ir.Read("password: ")
		conf.Email, _ = ir.Read("email: ")
		conf.Name, _ = ir.Read("name: ")
		conf.SmtpHost, _ = ir.Read("smtp host: ")
		for {
			smtp_port, _ := ir.Read("smtp port: ")
			port, err := strconv.Atoi(smtp_port)
			conf.SmtpPort = port
			if nil != err {
				log.Println(err)
				continue
			}
			break
		}
		fmt.Println("------------------")

		err = conf.Save(DEFAULT_CONFIG_FILE)
		return conf, err
	}
	return conf, err
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

	conf, err := getConfig(CONFIG_FILE)
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
