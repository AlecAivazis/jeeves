package main

// ReportError sends the error to the specified channel
func (b *JeevesBot) ReportError(channel string, errorToReport error) (err error) {
	// send the error message to the channel
	_, err = b.Discord.ChannelMessageSend(channel, errorToReport.Error())

	return err
}
