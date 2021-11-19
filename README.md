### 雪花算法
高并发环境下对唯一ID生成算法，算法来自 [https://github.com/asong2020/go-algorithm/tree/master/snowFlake](https://github.com/asong2020/go-algorithm/tree/master/snowFlake) 

对原来算法做了少量调整，把机器号调整为12位，其中8位为机器ID，4位数据中心ID，10位序列号，每毫秒最多能生成1024位序列号