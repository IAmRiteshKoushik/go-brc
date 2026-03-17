package main

func consumer(channel chan []byte) {
	for {
		<-channel
	}
}

