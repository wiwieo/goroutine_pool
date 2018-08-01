# goroutine_pool
协程池，用于管理协程

分成两种做法

一、被动等待（goroutin_pool.go）

worker循环从job中获取内容，worker会一直等待（此处做了超时机制，等待10秒还没有读取到，则退出）直到读取到内容后，开始执行


二、主动通知（goroutin_pool_new.go）

job中有数据后，通知一个worker去处理，处理完成之后，将worker放回池中。如果worker不足且没到上限，则新建一个。如果到达上限，则等待其他处理完之后。


