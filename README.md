# pingpong
Simple design ping pong pattern implemented in a different ways

The pattern is to exchange data between two competitive goroutines. The basic idea is not just to exchange, but to wait for the action to be performed, i.e., did it yourself, let someone else do it, and so on. 

The pattern is widely known, the only question is how to synchronize the actions. Below you can see how to synchronize the two goroutines using channels, mutexes, and atomics.

golang: 1.20.2
cpu: 11th Gen Intel(R) Core(TM) i5-1135G7 @ 2.40GHz


```
|                                Method |       |              |          |
|-------------------------------------- |------:|-------------:|---------:|
|            BenchmarkPingpong_lock-8   |  1512 | 784616 ns/op | 197 B/op |
|            BenchmarkPingpong_ch-8     |  1554 | 699865 ns/op | 307 B/op |
|            BenchmarkPingpong_atomic-8 |  7989 | 140181 ns/op | 125 B/op |
