package observer

type Event string

const (
	EventStudentRegistered Event = "student.registered"
	EventStudentEnrolled   Event = "student.enrolled"
	EventGradeAdded        Event = "grade.added"
)

type Observer interface {
	Notify(event Event, data interface{})
}

type Notifier struct {
	observers map[Event][]Observer
}

func NewNotifier() *Notifier {
	return &Notifier{
		observers: make(map[Event][]Observer),
	}
}

func (n *Notifier) Subscribe(event Event, observer Observer) {
	n.observers[event] = append(n.observers[event], observer)
}

func (n *Notifier) Notify(event Event, data interface{}) {
	for _, obs := range n.observers[event] {
		obs.Notify(event, data)
	}
}
