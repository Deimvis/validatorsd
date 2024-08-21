package hi

import "errors"

type MyStruct struct {
    Answer int
}

func (s *MyStruct) ValidateSelf() error {
    if s.Answer != 42 {
        return errors.New("answer is wrong")
    }
    return nil
}
