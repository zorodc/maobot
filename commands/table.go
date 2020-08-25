package commands

// Used for interraction
//import logs "github.com/zorodc/maobot/loggers"
import "github.com/zorodc/maobot/dynamic"

type Command struct {
	Function     dynamic.RtFunc
	Arity        uint
	OptionalArgs []interface{}
	Description  string
	Usage        string
}

var Table = map[string]Command{}

/*
import "os"
import "layeh.com/gumble/gumble"
import "layeh.com/gumble/gumbleffmpeg"

func Playtest(c *gumble.Client) {
	file, err := os.Open("/Shared/Zara.opus")
	if err != nil {
		panic(err);}
	source := gumbleffmpeg.SourceReader(file)
	stream := gumbleffmpeg.New(c, source)
	if stream == nil {
		logs.Log(logs.InterractionLogs, "Failed to open file.")
		return; }
	logs.Log(logs.InterractionLogs, "Playing...")
	err = stream.Play()
	if err != nil {
		panic(err);}
//	stream.Wait()
//	logs.Log(logs.InterractionLogs, "Done playing.")
}
*/
