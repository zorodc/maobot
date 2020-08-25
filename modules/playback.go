/* This source file describes the general faculty of playing music. */
package modules
import "github.com/zorodc/maobot/commands"

type Player interface {
	//	Stop()
	Next()
	Paused() bool
	Pause()
	Play()
	Info() string
	Volume() float32
	SetVolume(float32)
}

var gCurrentPlayer Player

func init() {
	// pop|skip|next
	commands.Table["next"]    =
		commands.Command{Function:func(_ interface{}) { gCurrentPlayer.Next(); }}
	commands.Table["skip"]    = commands.Table["next"]
	commands.Table["pop"]     = commands.Table["next"]

	commands.Table["pause"]   =
		commands.Command{Function:func(_ interface{}) { gCurrentPlayer.Pause(); }}
	commands.Table["unpause"] =
		commands.Command{Function:func(_ interface{}) { gCurrentPlayer.Play(); }}
	commands.Table["play"]    = commands.Table["unpause"]
	
	commands.Table["info"]    =
		commands.Command{Function:func(_ interface{}) { gCurrentPlayer.Info(); }}

	commands.Table["volume"]  = commands.Command{Function:
		func(_ interface{}, vol float32) {gCurrentPlayer.SetVolume(vol); }}
	commands.Table["volup"]   = commands.Command{Function:
		func(_ interface{}, vol float32) {
			gCurrentPlayer.SetVolume(gCurrentPlayer.Volume() + vol); }}
	commands.Table["voldown"] = commands.Command{Function:
		func(_ interface{}, vol float32) {
			gCurrentPlayer.SetVolume(gCurrentPlayer.Volume() - vol); }}
	commands.Table["volumeup"]   = commands.Table["volup"]
	commands.Table["volumedown"] = commands.Table["voldown"]
}

func SetPlayer(player Player) {
	if gCurrentPlayer != nil {
		gCurrentPlayer.Pause()
	} else if gCurrentPlayer != player { 
		gCurrentPlayer = player
	}
}
