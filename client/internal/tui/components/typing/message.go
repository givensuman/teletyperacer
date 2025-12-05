package typing

// TypingCompletedMsg is sent when the user completes typing the text
type TypingCompletedMsg struct{}

// AccuracyMsg contains the typing accuracy
type AccuracyMsg struct {
	Accuracy float64
}

// ProgressMsg contains the typing progress
type ProgressMsg struct {
	Progress float64
}
