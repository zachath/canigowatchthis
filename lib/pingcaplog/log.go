package pingcaplog

import (
	"fmt"

	"github.com/pingcap/errors"
)

func PingcapMarshalStack(err error) any {
	if errors.HasStack(err) {
		st := errors.GetStackTracer(err).StackTrace()

		state := &state{}
		out := make([]string, 0, len(st))
		for _, frame := range st {
			source := frameField(frame, state, 's')
			line := frameField(frame, state, 'd')
			funcName := frameField(frame, state, 'n')

			trace := fmt.Sprintf("%s:%s (%s)", source, line, funcName)

			out = append(out, trace)
		}

		return out[:len(out)-2]
	}

	return nil
}

type state struct {
	b []byte
}

// To live up to the fmt.State interface
func (s *state) Write(b []byte) (int, error) {
	s.b = b
	return len(b), nil
}

func (s *state) Width() (int, bool) {
	return 0, false
}

func (s *state) Precision() (int, bool) {
	return 0, false
}

func (s *state) Flag(c int) bool {
	return false
}

func frameField(f errors.Frame, s *state, c rune) string {
	f.Format(s, c)
	return string(s.b)
}
