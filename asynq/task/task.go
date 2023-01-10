package task

// TypeImageResize A list of task types.
const (
	TypeImageResize   = "image:resize"
	TypeArchiveCreate = "archive:create"
)

type ImageResizePayload struct {
	UserID   int
	ImageUrl string
}

type ArchiveCreatePayload struct {
	UserID  int
	SrcDir  string
	DestDir string
}

// Task represents a unit of work to be performed.
type Task struct {
	// typename indicates the type of task to be performed.
	typename string

	// payload holds data needed to perform the task.
	payload []byte
}

func (t *Task) Type() string    { return t.typename }
func (t *Task) Payload() []byte { return t.payload }

// NewTask creates a task with the given typename, payload and ResultWriter.
func NewTask(typename string, payload []byte) *Task {
	return &Task{
		typename: typename,
		payload:  payload,
	}
}
