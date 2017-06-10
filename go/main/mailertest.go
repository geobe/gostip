// [_Command-line flags_](http://en.wikipedia.org/wiki/Command-line_interface#Command-line_option)
// are a common way to specify options for command-line
// programs. For example, in `wc -l` the `-l` is a
// command-line flag.

package main

// Go provides a `flag` package supporting basic
// command-line flag parsing. We'll use this package to
// implement our example command-line program.
import (
	"flag"
	"github.com/SlyMarbo/gmail"
	"log"
)

func main() {

	account := flag.String("mail", "me@nohwere.wow", "a mail account")
	passwd := flag.String("pw", "foo", "password of the mail account")

	flag.Parse()

	email := gmail.Compose("Email subject", "Email body should be a nice html template")
	email.From = *account
	email.Password = *passwd

	// Defaults to "text/plain; charset=utf-8" if unset.
	email.ContentType = "text/html; charset=utf-8"

	// Normally you'll only need one of these, but I thought I'd show both.
	email.AddRecipient("georg.beier@posteo.de")

	err := email.Send()
	if err != nil {
		log.Fatal(err)
	}

}
