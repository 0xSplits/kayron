package image

import (
	"time"
)

// Cooler is configured to return a dynamically adjusted wait duration for this
// worker handler to sleep before running again. The introduced jitter has the
// purpose of spreading out the same type of work across time, so that we ease
// the load on our dependency APIs, here ECR, and effectively try to prevent
// rate limits. E.g. a jitter of 1% applied to 1h results in execution variation
// of +-36s.
func (h *Handler) Cooler() time.Duration {
	return h.jit.Percent(1 * time.Hour)
}
