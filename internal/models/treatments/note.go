package treatments

import (
	"fmt"
	"time"
)

type NoteType int

const (
	Default NoteType = iota
	Announcement
	Exercise
)

type NoteEvent struct {
	enteredBy string
	createdAt time.Time
	note      string
	noteType  NoteType
	duration  float32
}

func (n NoteEvent) TreatmentKind() TreatmentType {
	return Note
}

func (n NoteEvent) TreatmentTime() time.Time {
	return n.createdAt
}

func (n NoteEvent) TreatmentEnteredBy() string {
	return n.enteredBy
}

func (n NoteEvent) String() string {
	switch n.noteType {
	case Exercise:
		return fmt.Sprintf("%s Exercise: %s for %d min", n.createdAt.Format(time.RFC3339), n.note, int(n.duration))
	case Announcement:
		return fmt.Sprintf("%s Announcement: %s", n.createdAt.Format(time.RFC3339), n.note)
	case Default:
		fallthrough
	default:
		return fmt.Sprintf("%s Note: %s", n.createdAt.Format(time.RFC3339), n.note)
	}
}

func NewNoteEvent(createdAt time.Time, enteredBy string, note string, duration float32) NoteEvent {
	return NoteEvent{createdAt: createdAt, enteredBy: enteredBy, note: note, duration: duration}
}
