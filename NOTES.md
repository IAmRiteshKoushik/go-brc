This is very hardware dependent performance. If you are on better hardware, 
then you can extract better results. More cores, will allow you to utilize the 
CPU and the RAM better. If you have GPUs then you'll be much much faster.

## Attempt 1 - Brute force

Naively read everything and process them in a hashmap. Took 600+ seconds.

## Attempt 2 - Default Buffer

Default buffer with 4096 bytes = 4KB

## Attempt 3 - Scanner Buffer

1. Buffer `32 x 32 (1 KB)` - 31.7s
2. Buffer `64 x 64 (4 KB)` - 23.64s
3. Buffer `128 x 128 (16 KB)` - 22.6s 
4. Buffer `256 x 256 (64 KB)` - 22.15s - > optimal
5. Buffer `512 x 512 (256 KB)` - 27.07s
6. Buffer `1024 x 1024 (1 MB)` - 27.74s

## Attempt 4 - reader.ReadBytes()

Used ~48 seconds, so it's a degradation

## Attempt 5 - reader.ReadLine()

Used ~24 seconds with default buffer size of 4KB. Did not experiment with other 
buffer sizes.

## Attempt 6 - file.Read()

This one is giving good results.

1. Buffer `32 x 32 (1 KB)` - 21.37s
2. Buffer `64 x 64 (4 KB)` - 16.6s
3. Buffer `128 x 128 (16 KB)` - 15.5s
4. Buffer `256 x 256 (64 KB)` - 15.4s
5. Buffer `512 x 512 (256 KB)` - 12.5s -> optimal
6. Buffer `1024 x 1024 (1 MB)` - 15s

## Attempt 7 - file.Read() with single Goroutine

In this one, I am extracting all the data out through a channel. There is a 
slight spike in performance. My 12.5 seconds reading went upto 15.3s due to 
communication overhead.

For a single goroutine - 15.33s

## Attempt 8 - file.Read() with multile goroutines as consumers

The consumer might read inconsistent data. This is because slices are reference 
types and when you are sending a slice into a channel, you are not sending a copy 
of it but the header that points to the same memory address in RAM.

Here's the problem it causes:
`Producer`: Reads data from a file into a reusable buffer
`Producer` send the buffer into the channel
`Consumer` receives the buffer and starts reading from it
`Producer` simulatenously moves to the next iteration, flushes the buffer and 
reloads it with new file.Read(buffer)

The problem, is the consuemr is still trying to process some parts of the data, 
and the producer just overwrote the buffer with new data from the file.

This can be fixed through `copying`. It creates an snapshot of the data before 
transmission. The consumer gets a private copy of the info instead of a shared 
buffer.

This time, the average is coming out to be `17.3 seconds` across 4 tests. Copying 
indeed has a certain overhead to it.

## Attempt 9 - Actually useful consumers (leftover logic to process stuff)
