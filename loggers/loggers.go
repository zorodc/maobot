/* Defines an interface for producing logger objects that may 
   recieve log entries of various kinds, and providing log entries of
   particular kinds to the various loggers. */
package loggers

import "io"
import "log" // Use standard library loggers to implement WriterLogger
import "os"
import "fmt"
import "layeh.com/gumble/gumble"

type LogKind int;

const (
	ErrorLogs LogKind = iota // Failures
	DebugLogs                // Internal actions
	InterractionLogs         // Interraction with the user
	// StatusLogs?
)

// Global list of loggers used when logging.
var loggerList = map[LogKind]([]Logger){}

// Called to log an entry to all receptive loggers.
func Log(kind LogKind, entry string) {
	for _, logger := range loggerList[kind] {
		logger.Print(entry + "\n")
	}
}
func Logf(kind LogKind, format string, args ...interface{}) {
	Log(kind, fmt.Sprintf(format, args...))
}

// Called to log a gumble event in an interesting way.
func LogGumbleEvent(e interface{}) {
//	Log(fmt.Sprintf("Got event: %v", reflect.TypeOf(e)), DebugLogs)
	switch ev := e.(type) {
	case *gumble.ACLEvent:
		Logf(DebugLogs, "Got an ACL: %v.", ev)

	case *gumble.BanListEvent:
		Logf(DebugLogs, "Got a BanList: %v.", ev)

	case *gumble.ContextActionChangeEvent:
		Logf(DebugLogs, "Got a ContextActionChange: %v.", ev);

	case *gumble.PermissionDeniedEvent:
		Logf(DebugLogs, "Got a PermissionDeniedEvent: %v.", ev)

	case *gumble.ConnectEvent:
		Log(DebugLogs, "Connected to channel: {" + ev.Client.Self.Channel.Name +
    	     "} as {" + ev.Client.Self.Name + "}.")

	case *gumble.DisconnectEvent:
		var reason string
			switch {
			case len(ev.String) > 0:
				reason = ev.String
			case ev.Type == gumble.DisconnectError:
				reason = "unknown error"
			case ev.Type == gumble.DisconnectBanned:
				reason = "banned"
			case ev.Type == gumble.DisconnectKicked:
				reason = "kicked"
			case ev.Type == gumble.DisconnectUser:
				return // ignore
			}
			Log(ErrorLogs, "Disconnected from server: " + reason)

	case *gumble.ServerConfigEvent:
		if (ev.MaximumMessageLength != nil) {
			Logf(DebugLogs, "Maximum message length is %d.",
				*ev.MaximumMessageLength)
		}
		if (ev.WelcomeMessage != nil) {
			Log(DebugLogs, "Welcome message: " + *ev.WelcomeMessage + ".")
		}

	case *gumble.ChannelChangeEvent:
		Logf(DebugLogs, "Got a ChannelChange: %v.", ev);

	case *gumble.TextMessageEvent:
		var recip string
		var acc []string
		for _, userptr := range ev.Users {
			acc = append(acc, userptr.Name)
		}
		for _, channelptr := range ev.Channels {
			acc = append(acc, channelptr.Name)
		}
		if (len(ev.Users) > 0) {
			recip += fmt.Sprintf("users:%v",acc)
		}
		if (len(ev.Channels) > 0) {
			recip += fmt.Sprintf("channels:%v",acc)
		}

		sender := "<no sender>"
		if (ev.Sender != nil) {
			sender = ev.Sender.Name
		}

		Logf(DebugLogs, "{%s} => %s: {`%s`}", sender, recip, ev.Message)
			

	case *gumble.UserChangeEvent:
		var action string
		switch {
		case (ev.Type & gumble.UserChangeConnected) != 0:
			action = "connected to the server"
		case (ev.Type & gumble.UserChangeDisconnected) != 0:
			action = "disconnected"
		case (ev.Type & gumble.UserChangeKicked) != 0:
			action = "was kicked"
		case (ev.Type & gumble.UserChangeBanned) != 0:
			action = "was banned"
		case (ev.Type & gumble.UserChangeRegistered) != 0:
			action = "was registered"
		case (ev.Type & gumble.UserChangeUnregistered) != 0:
			action = "was unregistered"
		case (ev.Type & gumble.UserChangeName) != 0:
			action = "changed their name"
		case (ev.Type & gumble.UserChangeChannel) != 0:
			action = "moved to " + ev.User.Channel.Name
//    These are currently ignored.
//		case gumble.UserChangeComment, gumble.UserChangeAudio,
//			gumble.UserChangeTexture, gumble.UserChangePrioritySpeaker,
//			gumble.UserChangeRecording, gumble.UserChangeStats
		default:
			return
		}
	
		user := "<no user>"
		if (ev.User != nil) {
			user = ev.User.Name
		}
		
		Logf(DebugLogs, "UserEvent: {%s} %s.", user, action)

	case *gumble.UserListEvent:
// Lists all the registered users, once a new user is registered.
// Currently ignored.

	default:
		panic("Event handling incomplete!")
	}	
}

func AddLogger(l Logger, kinds ...LogKind) {
	for _, kind := range kinds {
		loggerList[kind] = append(loggerList[kind], l)
	}
}

type Logger interface {
	Print(string)
}

// Wrapper type for log.Logger
// The wrapper is necessary because logger.Print(interface{})
// doesn't satisfy an interface with the method Print(string).
// This would be a reasonable thing to do, but unfortunately, it is not done.
type WriterLogger struct {
	logger *log.Logger
}

func NewWriterLogger(w io.Writer) Logger {
	return WriterLogger{logger:log.New(w, "", log.Ldate|log.Ltime)}
}

func (this WriterLogger) Print(s string) {
	this.logger.Print(s)
}

type MessageLogger struct {
	client *gumble.Client // can be nil
	maxLen uint // maximum message size the server can recieve
}

func NewMessageLogger(c *gumble.Client) MessageLogger {
	return MessageLogger{client:c}
}

func (this *MessageLogger) SetClient(c *gumble.Client) {
	if c == nil { panic("Client should not be set to nil."); }
	(*this).client = c
}

func (this *MessageLogger) SetMaxMsgLen(len uint) {
	this.maxLen = len
}

func (this MessageLogger) Print(message string) {
	// Silently avoid printing a message to a null reciever.
	if this.client == nil { return; }
	// If maxLen is 0, the server sent no maximum length.
	if this.maxLen > 0 && uint(len(message)) > this.maxLen {
		fmt.Fprintf(os.Stderr,
			"ERROR: Message {`" + message + "`} exceeds maximum length " +
				string(this.maxLen) + ", trying to send anyways...")
	}
	this.client.Self.Channel.Send(message, false)
}
