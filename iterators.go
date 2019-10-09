package discordgo

import (
	"container/list"
	"sort"
)

// HistoryIterator contains the information needed to iterate
// over the history of a Messageable using HistoryIterator.Iter()
type HistoryIterator struct {
	channel      Messageable
	messages     *list.List
	limit        int
	before       string
	after        string
	around       string
	reverse      bool
	retrieveWith func(retrieve int) []*Message
}

// NewHistoryIterator creates a new HistoryIterator based on the given Messageable
// and using the defaults; limit = 100 and going from newest to oldest.
// channel : The Messageable to get the history from
func NewHistoryIterator(channel Messageable) *HistoryIterator {
	iterator := &HistoryIterator{
		channel:  channel,
		limit:    100,
		messages: list.New(),
	}
	iterator.retrieveWith = iterator.retrieveBefore
	return iterator
}

// SetLimit sets the limit for how much messages to fetch.
// limit : the limit to use
func (h *HistoryIterator) SetLimit(limit int) *HistoryIterator {
	h.limit = limit
	return h
}

// SetBefore sets the message ID from which to iterate back to get the messages before it
// beforeID : the Message ID
func (h *HistoryIterator) SetBefore(beforeID string) *HistoryIterator {
	h.before = beforeID
	h.after = ""
	h.around = ""
	h.retrieveWith = h.retrieveBefore
	return h
}

// SetAfter sets the message ID from which to iterate back from till the most recent or limit
// afterID : the Message ID
func (h *HistoryIterator) SetAfter(afterID string) *HistoryIterator {
	h.after = afterID
	h.before = ""
	h.around = ""
	h.retrieveWith = h.retrieveAfter
	return h
}

// SetAround sets the message ID from which to get the surrounding messages (max 100)
// aroundID : the Message ID
func (h *HistoryIterator) SetAround(aroundID string) *HistoryIterator {
	h.around = aroundID
	h.before = ""
	h.after = ""
	h.retrieveWith = h.retrieveAround
	return h
}

// Reverse toggles reverse, by default reverse is off and messages are returned newest first
func (h *HistoryIterator) Reverse() *HistoryIterator {
	h.reverse = !h.reverse
	return h
}

// Reset sets everything back to defaults
func (h *HistoryIterator) Reset() *HistoryIterator {
	h.after = ""
	h.limit = 100
	h.reverse = false
	h.before = ""
	h.around = ""
	h.messages.Init()
	h.retrieveWith = h.retrieveBefore
	return h
}

// Iter allows for the iterator to actually iterate and thus be used in a for ... range expression.
func (h *HistoryIterator) Iter() <-chan *Message {
	if h.around != "" && h.limit > 100 {
		h.limit = 100
	}
	if h.limit <= 0 {
		h.limit = 100
	}
	ch := make(chan *Message, 100)
	go func() {
		for {
			if h.messages.Len() == 0 {
				h.fetchNextMessages()
			}

			if h.messages.Len() == 0 {
				close(ch)
				return
			}

			el := h.messages.Front()
			ch <- h.messages.Remove(el).(*Message)
		}
	}()
	return ch
}

func (h *HistoryIterator) fetchNextMessages() {
	retrieve := h.limit
	if retrieve <= 0 {
		return
	} else if retrieve > 100 {
		retrieve = 100
	}

	messages := h.retrieveWith(retrieve)

	if h.reverse {
		toReverse := make(TimeSorter, 0, len(messages))
		for _, m := range messages {
			toReverse = append(toReverse, m)
		}

		sort.Sort(sort.Reverse(toReverse))

		messages = make([]*Message, 0, toReverse.Len())
		for _, m := range toReverse {
			messages = append(messages, m.(*Message))
		}
	}

	for _, m := range messages {
		h.messages.PushBack(m)
	}

	if len(messages) < 100 {
		h.limit = 0
	}
}

func (h *HistoryIterator) retrieveBefore(retrieve int) []*Message {
	messages, _ := h.channel.GetHistory(retrieve, h.before, "", "")
	if len(messages) != 0 {
		h.limit -= retrieve
		h.before = messages[len(messages)-1].ID
	}
	return messages
}

func (h *HistoryIterator) retrieveAfter(retrieve int) []*Message {
	messages, _ := h.channel.GetHistory(retrieve, "", h.after, "")
	if len(messages) != 0 {
		h.limit -= retrieve
		h.after = messages[0].ID
	}
	return messages
}

func (h *HistoryIterator) retrieveAround(retrieve int) []*Message {
	if h.around != "" {
		messages, _ := h.channel.GetHistory(retrieve, "", "", h.around)
		h.around = ""
		return messages
	}
	return nil
}
