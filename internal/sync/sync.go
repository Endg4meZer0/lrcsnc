package sync

func Start() {
	// Goroutine of the position synchronizer
	go positionSynchronizer()

	// Goroutine to watch for player controller signals
	go messageReceiver()

	// Goroutine to actively synchronize the lyrics with the song
	go lyricsSynchronizer()
}
