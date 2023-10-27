package env

import (
	"os"

	"github.com/Kagami/go-face"
)

// ================================================================
//
// ================================================================
type Face struct {
	*face.Recognizer
}

// ================================================================
//
// ================================================================
func (e *Face) Open() error {
	var err error
	e.Close()
	e.Recognizer, err = face.NewRecognizer(os.Getenv("DIR_FACE_RECOGNIZATION_MODELS"))
	return err
}

func (e *Face) Close() {
	if e.Recognizer != nil {
		e.Recognizer.Close()
	}
}