package main

import (
	"bufio"
	"fmt"
	"net"
)

// 定义协程池类型
type Pool struct {
	worker_num  int           // 协程池最大worker数量,限定Goroutine的个数
	JobsChannel chan net.Conn // 协程池内部的任务就绪队列
}

// 创建一个协程池
func NewPool(cap int) *Pool {
	p := Pool{
		worker_num:  cap,
		JobsChannel: make(chan net.Conn),
	}
	return &p
}

// worker处理函数
func process(conn net.Conn, work_ID int) {
	defer conn.Close() // 关闭连接
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		workerHead := fmt.Sprintf("[worker %d]", work_ID)
		n, err := reader.Read(buf[:]) // 读取数据
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println(workerHead + "收到client端发来的数据：" + recvStr)
		conn.Write([]byte(workerHead + " success")) // 发送数据
	}
}

// 协程池中每个worker的功能
func (p *Pool) worker(work_ID int) error {
	//worker不断的从JobsChannel内部任务队列中拿Conn
	for conn := range p.JobsChannel {
		//如果拿到Conn,则执行对应处理
		process(conn, work_ID)
	}
	return nil
}

// 协程池Pool开始工作
func (p *Pool) Run(listenAddr string) {
	// 设置监听端口
	listen, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Printf("server fail to listen")
		return
	}

	// 首先根据协程池的worker数量限定,开启固定数量的Worker,
	// 每一个Worker用一个Goroutine承载
	for i := 0; i < p.worker_num; i++ {
		go p.worker(i)
	}

	// 将新申请的连接加入到就绪队列
	for {
		conn, err := listen.Accept() // 建立连接
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		p.JobsChannel <- conn
	}
}

// 关闭协程池
func (p *Pool) Close() {
	close(p.JobsChannel)
}

func main() {
	// 创建一个协程池,最大开启3个协程worker
	p := NewPool(3)

	// 设定监听，启动协程池p
	p.Run("127.0.0.1:20000")

	// 关闭channel
	p.Close()

}
