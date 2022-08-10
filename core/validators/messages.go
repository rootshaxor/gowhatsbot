package validators

import (
	gerrors "main/core/errors"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// Validating the interface is an events.message and the sender is me.
func ValidateFromeMe(i interface{}, _ *whatsmeow.Client) (bool, error) {
	switch e := i.(type) {
	case *events.Message:
		if e.Info.IsFromMe {
			return false, nil
		} else {
			return false, gerrors.NewGoWhatsBot(2, "Validator", "Not from me")
		}
	default:
		return false, gerrors.NewGoWhatsBot(2, "Validator", "Not a message")
	}
}
