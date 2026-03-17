package main

// Simple consumer, does not process anything .. was being used just for testing
func consumer(channel chan []byte) {
	for {
		<-channel
	}
}
