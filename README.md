# 1 Billion Rows Challenge

An attempt at the 1 billion rows challenge in Golang. (inspired from Primeagen's 
article reads)

Goal : Read `measurements.txt` which contains 1 billion rows of data in the 
format of :
```bash
<weather-station>;<temperature>
```
And fine the min, max and average for each weather station.

This is an optimization problem where one is expected to bring down the timing 
required for reading the file and computing the data. For generating the data 
instructions are given in the main repository below.

For generating data, checkout [gunnarmorling/1brc](https://github.com/gunnarmorling/1brc).

## Summarization of My Attempts

> System Specs: AMD Ryzen5 3500U Radeon Vega Mobile Gfx 2.1GHz (16GB RAM) (no GPU)

**Attempt 01**
In the first attempt, I have tried a brute force approach. The code for this can 
be found in `attempt_1.go`. The program ran for 394.14 seconds.

**Attempt 02**
In this attempt, I have removed the code to actually compute any data and have 
just focused on testing our various approaches to optimize the file reading 
capacity. Have switched between the following :
1. Basic Scanner
2. Buffered Scanner
3. Bufio Reader
