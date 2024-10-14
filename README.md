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
