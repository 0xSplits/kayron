package image

import (
	"time"
)

func (h *Handler) Cooler() time.Duration {
	return h.jit.Percent(12 * time.Hour)
}
