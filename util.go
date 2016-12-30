package cron

// Start the cron scheduler in its own go-routine.
func (c *Cron) IsRunning() bool {
	return c.running
}

func (c *Cron) StopAndCleanAll() {
	c.Stop()
	c.entries = nil
}
