package main

// ReportError sends the error to the specified channel
func (b *JeevesBot) ReportError(channel string, err error) error {
	// send the error message to the channel
	_, err := b.Discord.ChannelMessageSend(channel, err.Error())
}
