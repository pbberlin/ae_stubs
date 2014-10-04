package main

import (
	"net/http"

	ae_mail "appengine/mail"
	go_mail "net/mail"

	"appengine"
	"appengine/memcache"
	"time"

	_ "strings"

	"bytes"
	"fmt"
)

const confirmMessage = `Thank you for creating an account!
Please confirm your email address by clicking on the link below:
%s`

/*
  https://developers.google.com/appengine/docs/python/mail/receivingmail

 	email-address:	 string@appid.appspotmail.com
 	is routed to
   /_ah/mail/string@appid.appspotmail.com

   peter@libertarian-islands.appspotmail.com
 	is routed to
   /_ah/mail/peter@libertarian-islands.appspotmail.com
*/
func emailReceive1(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	defer r.Body.Close()

	msg, err := go_mail.ReadMessage(r.Body)
	if err != nil {
		http.Error(w, "err is: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var from string

	if msg != nil {

		b1 := new(bytes.Buffer)

		// see http://golang.org/pkg/net/mail/#Message

		from = msg.Header.Get("from") + "\n"
		b1.WriteString("from: " + from)

		sbody := msg.Header.Get("subject") + "\n"
		b1.WriteString("body: " + sbody)

		when, _ := msg.Header.Date()
		swhen := when.Format("2006-01-02 - 15:04 \n")
		b1.WriteString(swhen)

		b1.ReadFrom(msg.Body)

		item := &memcache.Item{
			Key:        "latestEmail",
			Value:      b1.Bytes(),
			Expiration: 180 * time.Second,
		}

		if err := memcache.Set(c, item); err != nil {
			c.Errorf("error adding email to memcache: %v", err)
		} else {
			c.Infof("email successfully saved to memcache")
		}

	} else {
		c.Warningf("-empty msg- " + r.URL.Path)
	}

	var m map[string]string = nil
	m = make(map[string]string)
	m["sender"] = from
	emailSend(w, r, m)
}

func emailView(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	// Get the item from the memcache
	if item, err := memcache.Get(c, "latestEmail"); err == memcache.ErrCacheMiss {
		w.Write([]byte("item not in the cache"))
	} else if err != nil {
		s2 := fmt.Sprintf("error getting item: %v", err)
		w.Write([]byte(s2))
	} else {
		s := string(item.Value)

		s2 := fmt.Sprintf("the email is\n%s", s)
		w.Write([]byte(s2))
	}

}

func emailReceive2(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	defer r.Body.Close()

	/* alternative code from https://developers.google.com/appengine/docs/go/mail/
	not net/mail ,but
	*/
	var b2 bytes.Buffer
	if _, err := b2.ReadFrom(r.Body); err != nil {
		c.Errorf("Error reading body: %v", err)
		return
	}
	c.Infof("\n\nb2: " + b2.String() + "--\n ")

}

func emailSend(w http.ResponseWriter, r *http.Request, m map[string]string) {

	c := appengine.NewContext(r)
	//addr := r.FormValue("email")

	addr := m["sender"]
	_ = addr
	email_thread_id := []string{"3223"}

	msg := &ae_mail.Message{
		//Sender:  "Peter Buchmann <peter.buchmann@web.de",
		Sender: "peter.buchmann@web.de",
		//To:	   []string{addr},
		To: []string{"peter.buchmann@web.de"},

		Subject: "Confirm your registration",
		Body:    fmt.Sprintf(confirmMessage, "http://some_url"),
		Headers: go_mail.Header{"References": email_thread_id},
	}
	if err := ae_mail.Send(c, msg); err != nil {
		c.Errorf("Couldn't send email: %v", err)
	} else {
		c.Infof("email successfully sent")
	}

}
