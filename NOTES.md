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

## Attempt 9 - Actually useful buffers (leftover logic to process stuff)

In this iteration, I was fixing the issue that the size of the buffer is fixed 
but then, the problem comes when lines transported to the consumer might be 
cut from the middle. 

In order to handle this, I am going to only send complete lines by checking for 
newline characters from the end. After checking for them, I will only send the 
part of the buffer that's fully valid and keep the leftover in a separate 
buffer that will be pre-pend to the next batch of data being sent to the channel.

In this approach, copying seems to introduce minimal overhead and my timing is 
still around `17.4 seconds`

## Attempt 10 - Evolving the Consumer to make it actually useful

The grid test was failing initially due to some bugs in the code. Mostly it is 
a circular wait problem where I have goroutines waiting to read from a channel
but there is a wait group that is not letting them read.

Upon fixing that, and setting up item sizes with parallel works, the following 
is the grid test

Workers/Items   5    10   20   50   100
---------------------------------------
5               165  161  149  142  145
10              151  100  97   83   83
20              102  101  96   93   77
50              108  103  106  105  108

So, it's obvious that the best configuration is coming in 10 x 100 and 20 x 100
Going forward, these two are the only configurations to be tested.

## Attempt 11 - Custom Byte Split

A custom byte split function lead to the creation of `consumerV3`. This brought
down the timing from 77s. Although, the behavior is slightly erratic, here are 
the obeservations

Iteration 1: 50.56s
Iteration 2: 50.33s
Iteration 3: 52.00s
Iteration 4: 57.90s
Iteration 5: 56.68s
Iteration 6: 59.14s
Iteration 7: 54.58s
Iteration 8: 59.09s

I have come to realize that my CPU is heating up and a thermal throttling is 
kicking in. Between iterations 6 and 7, I waited a little and I could see a drop 
in the timing again.

So, I would like to consider the median value in this observation here that is 
55 seconds.

---

The justification for this performance gain is that now I am using a custom 
byte splitting function. This is allowing me to re-use a pre-defined byte 
buffer and save up on needless memory allocations.

## Attempt 12 - Custom Byte Hash using FNV Hash to remove string allocation

In this attempt, I am going to reduce the amount of data that I store in the 
hashmap. This time instead of using a string as a key, I am going to use a uint64
via FNV hashing. his is because, I can hash the bytes and generate this pretty easily.

Now, coming to the benefits of it: 
1. `map[string]any` -> 16 bytes of key space + N variable bytes based on size
    This is because a string consists of the actual data and a StringHeader
    The StringHeader contains a uintptr (8 bytes) and len (8 bytes)
2. `map[uint64]any` -> 8 bytes of key space (fixed)

In addition to this, lookups on the hashmap become faster. This is due to the 
fact that I am using a 64-bit architecture system. So, when I have to has the 
key to look it up in the hashmap, it fits within one register and the hashing 
is done in a single cycle.

Internally, a hashmap in Go uses buckets and each bucket contains 8 entries. 
This is to handle collision. I am not going into depth of why this is done 
the way it is.

However, what is interesting to note is that hashing is only one step of the 
lookup and this locates the apt bucket. Now, you need to scan each element of the 
bucket to determine a match.

This is done using the `XOR` operator. If fully matched then 0 is returned. Now, 
because the entire key fits into a single register, it takes the CPU a single 
cycle to compare the lookup-key against the 8 candidates of the bucket.

Results:
Iteration 1: 38.96s
Iteration 2: 40.12s
Iteration 3: 38.18s
Iteration 4: 38.64s

So the average time comes out to 38.8s

## Attempt 13 - Custom Float Parser

Every temperature is a single decimal float number so there is no need to use 
ParseFloat() which runs multiple sanity checks and does error handling. Is it 
simple to just read the 4 characters, skip the decimal and then return back the 
three characters in the form of an integer. In the final display, a decimal 
point can be added.

The mean timing came down to 31.5s

## Attempt 14 - Custom Scanner

Next up, according to the flamegraph, a lot of time gets wasted in the default 
scanner being used. This is the next site of optimization.

The scanner has mostly become redundant and its existance is causing unnecessary 
allocation on the heap. Also, the allocation happens for every line on input.

Here, the implementation remains straight forward. The idea is inspired from 
the parseLine() function. Previously, we were working with \n and then splitting 
around ;. Now, we are doing the same but also skipping the decimal point and 
returning the index of the next spot to start reading

The GC pressure is much lower now as we are managing the reading index ourselves 
as it is just an integer. Prevously the scanner was maintaining internal state  
on the heap that had to be deallocated, now we are just overwriting.

The observations are as follows: 27.12, 27.14, 27.97, 27.50. So we can say that 
27.5 is a good estimate to play with.

## Attempt 15 - Reducing Channel Communication Overhead with Worker Reading

## Attempt 16 - Delete the Name and Temperature Buffers
