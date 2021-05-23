package pool


type PooledObject interface{
	Create() interface{}
	Destroy()

}

//Pool Interface
type Pool interface{
	Capacity() int
	CurrentSize() int
	InitSize() int
	MinSize() int
	MaxSize() int
	IdleTimeout() int
	WaitTimeout() int
	HighWaterMark() int
	LowWaterMark() int
	Checkout() interface{}
	CheckIn() interface{}
	CheckOutCount() int64
	CheckInCount() int64
	CreateCount() int64
	DestroyCount() int64
	Start()
	Stop()
	Configure(maxPoolSize,minPoolSize,initialSize,idleTimeout,maxReusableCount,waitTime,retryCount int)

}

//PoolConfig struct holds the configuration of the pool
type PoolConfig struct{
	MaxPoolSize int
	MinPoolSize int
	InitialSize int
	IdleTimeout int
	MaxReusableCount int
	WaitTime int
	RetryCount int

}