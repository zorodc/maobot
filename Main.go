// This is a music-playing bot for the voice-chat program mumble.
// It has a number of commands, for playing, queueing, and pausing music.
// It uses youtube-dl as a backend, so it can play nearly anything.
// It's also a stoolie, and outputs the terminal what your users say.
// One final capability is that it can output images from links users send.
//
// Thus, if you send a link to an image, maobot will automatically
//  download and send the image itself to the chat.
//
// A tip - if you wish to see all commands it can be given,
//  I suggest grep -R'ing for 'commands.Table,' in this directory.

package main

import "layeh.com/gumble/gumble"
import gutil "layeh.com/gumble/gumbleutil"

import "github.com/zorodc/maobot/commands"
import logs "github.com/zorodc/maobot/loggers"
import "github.com/zorodc/maobot/eventstream"

//import "flag"
import "fmt"
import "os"
import "net"
import "crypto/tls"
import "time"
import "errors"
import "strings"

import "unicode"
import "unicode/utf8"

import _ "layeh.com/gumble/opus"
import _ "github.com/zorodc/maobot/modules"
import "github.com/zorodc/maobot/imgfetch"

const (
	// Maximum reconnection tries. 
	kDMaxTries = uint(8)
	kDTimeout = 4*time.Second
	kDPort    = "64738"
)

const Usage =
	`maobot: username[:password] address[:port] [-insecure] 
`
// Mandatory parameters returned, optional parameters taken as pointers.
func ParseArgs(defaultPort string, confarg *gumble.Config,
	dialer *net.Dialer, tlsConf *tls.Config,
	file **os.File) (addr string, conf *gumble.Config, err error) {
	conf = confarg // Allow the client to provide the config object.

	pullPair := func(s string, opt *string) string {
		idx := strings.LastIndex(s, ":")
		if idx == -1 {
			idx = len(s)
		} else {
			*opt = s[idx+1:]
		}
		return s[:idx]
	}

	if len(os.Args) < 3 {
		err = errors.New("Not enough arguments.")
		return
	}

	conf.Username = pullPair(os.Args[1], &conf.Password)
	url := pullPair(os.Args[2], &defaultPort)
	addr = url + ":" + defaultPort

	return
}

func skipWhiteSpace(msg string) string {
	for len(msg) > 0 {
		r, size := utf8.DecodeRuneInString(msg)
		if !unicode.IsSpace(r) { break; }
		msg = msg[size:]
	}
	return msg
}

// Attempt connection/reconnection `tries` number of times.
func TryReconn(tries uint,
	dialer *net.Dialer, tlsConf *tls.Config, address string,
	config *gumble.Config) (client *gumble.Client, err error) {
	for i := uint(0); i < tries; i++ {
		logs.Logf(logs.ErrorLogs, "Reconnecting... Try: %d", i+1)
		client, err = gumble.DialWithDialer(dialer, address, config, tlsConf)

		if err == nil { break; }
		time.Sleep(300 * time.Millisecond)
	}; return
}

func main() { /* youtube-dl -x -f opus/bestaudio -o - '...' */
	/* Persistent objects. */
	// MessageLogger: configured w/ setters once a connection is established.
	messagelogger := logs.NewMessageLogger(nil)
	conf          := gumble.NewConfig()
	dialer        := net.Dialer{Timeout: kDTimeout}
	tlsConf       := tls.Config{InsecureSkipVerify:true}
	address       := ""
	sigExit       := make(chan int)

	/* Parse arguments. */
	debugOut           := os.Stdout // default
	address, conf, err := ParseArgs(kDPort, conf, &dialer, &tlsConf, &debugOut)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Print(Usage)
		os.Exit(1);}

	/* Setup loggers. */
	logs.AddLogger(logs.NewWriterLogger(os.Stdout), logs.InterractionLogs)
	logs.AddLogger(logs.NewWriterLogger(debugOut),  logs.DebugLogs)
	logs.AddLogger(logs.NewWriterLogger(os.Stderr), logs.ErrorLogs)
	logs.AddLogger(&messagelogger, logs.ErrorLogs, logs.InterractionLogs)

	/* Attach event listeners. */
	// Main listener
	conf.Attach(gutil.ListenerFunc(func(e interface{}) {
		logs.LogGumbleEvent(e)
		eventstream.PostEvent(e)
		
		switch e := e.(type) {
		case *gumble.ServerConfigEvent:
			if (e.MaximumMessageLength != nil) {
				messagelogger.SetMaxMsgLen(uint(*e.MaximumMessageLength)); }

		case *gumble.TextMessageEvent:
			// Skip the leading whitespace that mobile clients add.
			plaintext := skipWhiteSpace(gutil.PlainText(&e.TextMessage))
			commands.Dispatch(plaintext, e.Client)

			// Fetch images from plain links sent by mobile users.
			b64, mt, ok := imgfetch.FetchImage(plaintext)
			if ok == nil {
				e.Client.Self.Channel.Send(
					`<img src="data:`+mt+`;base64,`+string(b64)+`"/>`, false)
			}

		case *gumble.ConnectEvent:
			messagelogger.SetClient(e.Client)
		}}))

	// Disconnect listener.
	conf.Attach(gutil.Listener{Disconnect:func(e *gumble.DisconnectEvent) {
		switch e.Type {
		case gumble.DisconnectError:
			// Gumble sometimes sends a DisconnectError while the connction is on.
			// So, we disconnect manually to "fix" this.
			e.Client.Disconnect()
			_, err := TryReconn(kDMaxTries, &dialer, &tlsConf, address, conf)
			if err != nil {
				logs.Logf(logs.ErrorLogs, "Couldn't reconnect after %d tries: %s.",
					kDMaxTries, err.Error())
				sigExit <- 1
			} else { logs.Logf(logs.DebugLogs, "Reconnection successful."); }

		case gumble.DisconnectKicked, gumble.DisconnectBanned:
			// Exit with failure code on kick or ban.
			sigExit <- 1
		case gumble.DisconnectUser:
			// Exit with success on voluntary disconnection.
			sigExit <- 0
		}}})

	/* Setup connection. */
	_, err = gumble.DialWithDialer(&dialer, address, conf, &tlsConf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failure to connect: ", err)
		os.Exit(1)
	}
	
	// Wait on an exit signal.
	os.Exit(<-sigExit)
}
