package eventstream

type EventRecipient func(interface{})bool

var recipients []EventRecipient

// Signal to all recipients that an event has occurred.
func PostEvent(event interface{}) {
	var done []uint
	for idx, recipient := range recipients {
		// Create a list of recipients to remove.
		if recipient(event) {
			done = append(done, uint(idx))
		}
	}
	var nremoved uint;
	for _, idx := range done {
		// Adjust for the recipients already removed from the slice.
		// e.g. If we removed one earlier in the list, the next removal needs to
		// adjust by 1.
		idx -= nremoved
		// Remove a single element.
		recipients = append(recipients[:idx], recipients[idx+1:]...)
		nremoved++
	}
	
}

func PostRecipient(er EventRecipient) {
	recipients = append(recipients, er)
}
