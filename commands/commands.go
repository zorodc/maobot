/* Defines an interface for dynamic function calls from strings, allowing for
   runtime function dispatch. Reflection is used internally, but clients are
   expected to interract using interface{} objects. */

// TODO
// And have the ServerCommands package have a sub-bpackage called "modules,"
// And have each module included into the ServerCommands main file and
// include itself into the table.
package commands

import logs "github.com/zorodc/maobot/loggers"
import      "github.com/zorodc/maobot/dynamic"

import "strings"
import "bytes"

/*	"ytsearch":todo,
	// pandora commands
	"pandora":todo,
	"channels":todo,
	"setchannel":todo,
	// heteroqueue commands
	"queue":todo,
	"add":todo, "append":todo,
	"playnext":todo, "prepend":todo,
	"playnow":todo,
	"prev":todo, "previous":todo, "prior":todo,
	"push":todo, "interrupt":todo,
	"clear":todo,
	"forward":todo,
	"back":todo,
	"scroll":todo,
	"dl_status":todo,
	// shared commands
	"stop":todo,
	"pop":todo, "skip":todo, "next":todo,
	"pause":todo,
	"unpause":todo, "play":todo,
	"info":todo,
	"volup":todo, "volumeup":todo,
	"voldown":todo, "volumedown":todo,
	"vol":todo, "volume":todo,
	// movement commands
	"moveto":todo,
	"summon":todo,
	"follow":todo, "harass":todo,
	// timing commands
	"remindme":todo,
	// scripting commands
	"script":todo,
	"batch":todo,
	// logging commands
	"pt_message":todo,
	"message":todo,
	"debugon":todo,
	"debugoff":todo,
	// miscellaneous commands
	"fortune":todo,*/

// Parses from a start quote to an end quote, allowing escaping of the " and \.
func parseQuotedString(quoted string) (uint, string) {
	var buffer bytes.Buffer
	var sawslash bool
	var nread uint
	for _, rune := range quoted {
		nread++
		switch {
		case !sawslash && rune == '\\': sawslash = true
		case !sawslash && rune == '"' : break
		default:
			sawslash = false
			buffer.WriteRune(rune)
		}
	}
	return nread, buffer.String()
}

// Translate the strings into instances of some type.
func evalType(ss []string) []interface{}{
	is := make([]interface{}, len(ss))
	for i := range ss {
		is[i] = ss[i]
	}; return is
}

func ParseArguments(msg string) (arguments []string) {
	for len(msg) > 0 {
		if msg[0] == '"' {
			nread, parsed := parseQuotedString(msg)
			arguments = append(arguments, parsed)
			msg = msg[nread:]
		} else if idx := strings.IndexRune(msg, ' '); idx != -1 {
			arguments = append(arguments, msg[:idx])
			// Skip spaces
			var i uint
			for i = uint(idx); i < uint(len(msg)) && msg[i] == ' '; i++ {}
			msg = msg[i:]
		} else {
			arguments = append(arguments, msg)
			break
		}
	}
	return
}

func Dispatch(msg string, param interface{}) {
	if !strings.HasPrefix(msg, "!") { return; } 
	lst := ParseArguments(msg)
	logs.Logf(logs.DebugLogs, "Attempting call to: %#v", lst)

	cmd  := lst[0]
	args := append([]interface{}{param}, evalType(lst[1:])...)
	logs.Logf(logs.DebugLogs, "ARGS:%#v", args)

	if command, ok := Table[cmd[1:]]; ok {
		_, err := dynamic.Call(command.Function, args...)
		if err != nil {
			logs.Logf(logs.ErrorLogs, "Improper call: %s.", err.Error())
			logs.Log(logs.InterractionLogs, "Usage: " + cmd + " " + command.Usage)
		}
	} else {
		logs.Logf(logs.DebugLogs,
			"Nonexistent command `%s` called with arguments %#v.", cmd, args)
		logs.Logf(logs.InterractionLogs, "No such command `%s`.", cmd)
		// logs.Log(CommandList(),logs.InterractionLogs)
	}
}

/*var todo dynamic.RtFunc = dynamic.New(todo_)
func todo_(e interface{}) {
	_, err := dynamic.Call(fmt.Printf, []interface{}{"%s %d %d", "hi", 13, 100}...)
	if err != nil {
		fmt.Printf("Error in call: %s.\n", err.Error())
	}
	_, err = dynamic.Call(logs.Log, []interface{}{logs.InterractionLogs, "hi!"}...)
	if err != nil {
		fmt.Printf("Error in call: %s.\n", err.Error())
	}

	_, err = dynamic.Call(logs.Log, logs.InterractionLogs, "hello from the dynamic world!")
	if err != nil {
		fmt.Printf("Error in call: %s.\n", err.Error())
	}
}
*/
