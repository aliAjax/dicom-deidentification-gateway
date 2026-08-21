package application

// Terminal states used by Recover to roll a spool item back to a safe,
// reprocessable state. These are defined in one place so the application
// layer owns the spool's state vocabulary rather than scattering string
// literals through the call sites.
const (
	// failedState is terminal: the item's spool file could not be read, so
	// retrying would only hit the same error. An operator must intervene.
	failedState = "failed"
	// recoveryFailureState leaves the item pending so a later recovery pass
	// can retry it once the transient handler error has cleared.
	recoveryFailureState = "pending"
	// doneState marks a fully recovered item so it is not reprocessed.
	doneState = "done"
)
