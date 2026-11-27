package ui

// quiet drops header, step, span and summary output but keeps results and
// errors, plus any Detail/Table that directly follows a kept call.
type quiet struct {
	next     Reporter
	lastKept bool
}

// Quiet wraps r so only results and errors are shown.
func Quiet(r Reporter) Reporter { return &quiet{next: r} }

func (q *quiet) Header(string, ...Field) { q.lastKept = false }
func (q *quiet) Step(Status, string)     { q.lastKept = false }
func (q *quiet) Summary(Summary)         { q.lastKept = false }

type noopSpan struct{}

func (noopSpan) Done(Status, string) {}

func (q *quiet) Start(string) Span {
	q.lastKept = false
	return noopSpan{}
}

func (q *quiet) Detail(lines ...string) {
	if q.lastKept {
		q.next.Detail(lines...)
	}
}

func (q *quiet) Table(h []string, rows [][]string) {
	if q.lastKept {
		q.next.Table(h, rows)
	}
}

func (q *quiet) Result(s Status, text string) {
	q.lastKept = true
	q.next.Result(s, text)
}

func (q *quiet) Error(err error, hint string) {
	q.lastKept = true
	q.next.Error(err, hint)
}
