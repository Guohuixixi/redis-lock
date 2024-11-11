# redis-lock
-----------------
基于redis setnx实现的分布式锁
+ 续约
  + 手动续约
  + 自动续约
+ 超时重试+重试策略
+ 利用singleflight优化分布式锁


