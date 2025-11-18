package button

// FocusMsg describes a tea.Msg on a button's
// focus state
type FocusMsg bool

const (
	Focus   FocusMsg = true
	Unfocus FocusMsg = false
)


// DisableMsg describes a tea.Msg on a button's
// disabled state
type DisableMsg bool

const (
	Disable DisableMsg = true
	Enable  DisableMsg = false
)

// WidthMsg is used to set the width of a button
type WidthMsg int
