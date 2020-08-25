// TODO: This ought to go in package modules/playback/heteroqueue?
// In the source file "player.go"?
// With a source file called "heteroqueue.go" containing add. functionality?
// With another source file called "commands.go" containing commands?

package modules

import "layeh.com/gumble/gumble"
import "layeh.com/gumble/gumbleffmpeg"
import "github.com/zorodc/maobot/commands"
import "github.com/zorodc/maobot/collections/syncqueue"
import "time"
//import "github.com/zorodc/maobot/eventstream"

const kLoopInterval = 200 * time.Millisecond

// TODO: On rollback, grab source of stopped song and restart w/ it.
// Otherwise unpause.

// Spinlock. Once the current song is done, go to the next one.
// TODO: Add synchronization to StreamQueue functions.
// This will ensure exclusivity of queue manipulation.
func thread() { 
	for {
		time.Sleep(kLoopInterval)

		s := gStreamQueue.Front()
		if s != nil &&
			s.(*gumbleffmpeg.Stream).State() == gumbleffmpeg.StateStopped {
			gStreamQueue.PopFront()
			if s := gStreamQueue.Front(); s != nil {
				s := s.(*gumbleffmpeg.Stream)
				err := s.Play()
				println("Skipping!")
				if err != nil {
					println(err.Error())
				}
			}
// TODO: FIX THIS FOR OTHERS
//			if s := gStreamQueue.Front().(*gumbleffmpeg.Stream); s != nil {
//				s.Play()
//			}
		}
	}
}

func init() {
//	eventstream.PostRecipient(func(e interface{}) bool {
//		switch e := e.(type) {
//		case *gumble.UserChangeEvent:
//			if (e.Type & gumble.UserChangeConnected) == 0 {
//				return false
//			}
//
//			println("Recieved event.")
//			stream := gumbleffmpeg.New(e.Client, source)
//			gStreamQueue.Append(stream)
//			gStreamQueue.Front().(*gumbleffmpeg.Stream).Play()
//		}
//		return false
//	})

	SetPlayer(&gStreamQueue)
	go thread()
	commands.Table["queue"] = commands.Command{
		Function:func(_ interface{}) { SetPlayer(&gStreamQueue); },
		Arity:0,
		OptionalArgs:nil,
		Description:"Set the queue as the current player.",
		Usage:"",}
	commands.Table["add"] = commands.Command{
		Function:func(c *gumble.Client, link string) {
			source := gumbleffmpeg.SourceExec(
				"youtube-dl", "-f", "opus/bestaudio", "-o", "-", link)
			println("Source object made.")
			stream := gumbleffmpeg.New(c, source)
			println("Stream object made.")
			gStreamQueue.Append(stream)
			println("Stream object appended.")
			err := gStreamQueue.Front().(*gumbleffmpeg.Stream).Play()
			println("Playing...")
			if err != nil {
				println(err.Error())
			}
		},
		Arity:1,
		OptionalArgs:nil,
		Description:"Set the queue as the current player.",
		Usage:"",}
/*	commands.Table["fromfile"] = commands.Command{
		Function:fromFile,
		Arity: 1,
		OptionalArgs:nil,
		Description:"Add a piece to the queue from a file.",
		Usage:"filename",}*/
	/*
		Function:playreplace,
		Arity:1,
		OptionalArgs:nil,
		Description:"Play an arbitrary song from the disk.",
		Usage:"filename",}*/
/*	commands.Table["stop"] = commands.Command{
		Function:func(_ interface{}) { stop(); },
		Arity:0,
		OptionalArgs:nil,
		Description:"Stop the currently playing song.",
		Usage:"",}*/
/*	commands.Table["toggle"] = commands.Command{
		Function:toggle,
		Arity:0,
		OptionalArgs:nil,
		Description:"Pause/unpause the currently playing song.",
		Usage:"",}*/
}

/* Implementation of the player interface for ffmpeg streams. */
type StreamPlayer struct {
	syncqueue.Queue
}

var gStreamQueue StreamPlayer

func (this *StreamPlayer) Next() {
	front := this.PopFront()
	if front != nil {
		front.(*gumbleffmpeg.Stream).Pause()
	}
	this.Play()
}

func (this *StreamPlayer) Paused() bool {
	front := this.Front()
	if front != nil {
		return front.(*gumbleffmpeg.Stream).State() == gumbleffmpeg.StatePaused
	}

	return false
}

func (this *StreamPlayer) Pause() {
	front := this.Front()
	if front != nil {
		front.(*gumbleffmpeg.Stream).Pause()
	}
}

func (this *StreamPlayer) Play() {
	front := this.Front()
	if front != nil {
		front.(*gumbleffmpeg.Stream).Play()
	}
}

func (this *StreamPlayer) Info() string {
	return "" // TODO
}

func (this *StreamPlayer) Volume() float32 {
	return 0.0 // TODO
}

func (this *StreamPlayer) SetVolume(vol float32) {
	println(vol) // TODO
}
