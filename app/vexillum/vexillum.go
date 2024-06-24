package vexillum

import caesar "github.com/caesar-rocks/core"

// Vexillum is a (very) lightweight feature flagging feature.
type Vexillum struct {
	Flags map[string]bool
}

func New() *Vexillum {
	return &Vexillum{
		Flags: make(map[string]bool),
	}
}

// IsActive returns true if the feature is active.
func (v *Vexillum) IsActive(feature string) bool {
	return v.Flags[feature]
}

// Activate activates a feature.
func (v *Vexillum) Activate(feature string) {
	v.Flags[feature] = true
}

// Deactivate deactivates a feature.
func (v *Vexillum) Deactivate(feature string) {
	v.Flags[feature] = false
}

func (v *Vexillum) Middleware(feature string) caesar.Handler {
	return func(ctx *caesar.Context) error {
		if v.IsActive(feature) {
			ctx.Next()
			return nil
		}
		return caesar.NewError(400)
	}
}
