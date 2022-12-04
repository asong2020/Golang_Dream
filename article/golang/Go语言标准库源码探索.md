## 前言

哈喽，大家后，我是`asong`；这几天看了一下Go语言标准库`net/http`的源码，所以就来分享一下我的学习心得；为什么会突然想看http标准库呢？因为在面试的时候面试官问我你知道Go语言的`net/http`库吗？他有什么有缺点吗？因为我没有看过这部分源码，所以一首凉凉送给我；

废话不多说，接下请跟着我的脚步我们一起探索`net/http`；

**本文代码基于：Go1.19.3**



## net/http库的一个小demo

服务端：

```go
import (
	"fmt"
	"net/http"
)

func getProfile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "asong")
}

func main() {
	http.HandleFunc("/profile", getProfile)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("http server failed, err: %v\n", err)
		return
	}
}
```

本地启动一个server端监听`8080`端口，并且提供路由`/profile`获取个人信息；

客户端：

```go
import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	resp, err := http.DefaultClient.Get("http://127.0.0.1:8080/profile")
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read from resp.body failed, err:%v\n", err)
		return
	}
	fmt.Println(string(body))
}
```

通过这样一个简单的例子，我们可以知道客户端我们主要使用`http.Client{}`，服务端我们主要使用`http.ListenAndServe`和`http.HandleFunc`，所以我们就可以从这两个包入手来分别看一看客户端和服务端的代码具体是怎么封装；



## 客户端实现

客户端级别最高的抽象是`net/http.Client{}`，具体结构如下：

```go
type Client struct {
    Transport RoundTripper
    CheckRedirect func(req *Request, via []*Request) error
    Jar CookieJar
    Timeout time.Duration
}
```

- Transport：其类型是RoundTripper，RoundTrip代表一个本地事务，RoundTripper接口的实现主要有三个：Transport、http2Transport、fileTransport，其目的是支持更好的扩展性；
- CheckRedirect：用来做重定向
- Jar：其类型是CookieJar，用来做cookie管理，CookieJar接口的实现Jar结构体在源码包`net/http/cookiejar/jar.go`；

客户端可以直接通过`net/http.DefaultClient`发起HTTP请求，也可以自己构建新的`net/http.Client`实现自定义的HTTP事务，多数情况下我们使用默认的客户端发出的请求就可以满足需求；

我们画一个UML图，如下所示：

![image-20221204160448320](C:\Users\sunsong\AppData\Roaming\Typora\typora-user-images\image-20221204160448320.png)

了解HTTP客户端的基本结构，我们接下来就开始分析客户端的基本实现；



### 构建Request

`net/http`包的`Request`结构体封装好了HTTP请求所需的必要信息：

```go
type Request struct {
	Method string
	URL *url.URL
	Proto      string // "HTTP/1.0"
	ProtoMajor int    // 1
	ProtoMinor int    // 0
	Header Header
	Body io.ReadCloser
	GetBody func() (io.ReadCloser, error)
	ContentLength int64
	removed as necessary when sending and
	TransferEncoding []string
	Close bool
	Host string
	Form url.Values
	PostForm url.Values
	MultipartForm *multipart.Form
	Trailer Header
	RemoteAddr string
	RequestURI string
	TLS *tls.ConnectionState
	Cancel <-chan struct{}
	Response *Response
	ctx context.Context
}
```

其中包含了HTTP请求的方法、URL、协议版本、协议头以及请求体等字段，还包括了指向响应的引用：Response；其提供了`NewRequest()、NewRequestWithContext()`两个方法用来构建请求，这个方法可以校验HTTP请求的字段并根据输入的参数拼装成新的请求结构体，`NewRequest()`方法内部也是调用的`NewRequestWithContext`，区别就是是否使用context来做goroutine上下文传递；接下来我们看一下`NewRequestWithContext`方法的具体实现：

```go
func NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*Request, error) {
    // 默认使用GET方法
	if method == "" {
		method = "GET"
	}
    // 校验方法是否有效，也就是是否是GET、POST、PUT等
	if !validMethod(method) {
		return nil, fmt.Errorf("net/http: invalid method %q", method)
	}
    // ctx必须要传递，NewRequest方法调用时会传递context.Background()
	if ctx == nil {
		return nil, errors.New("net/http: nil Context")
	}
    // 解析URL，解析Scheme、Host、Path等信息
	u, err := urlpkg.Parse(url)
	if err != nil {
		return nil, err
	}
    // body在下面会根据其类型包装成io.ReadCloser类型
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = io.NopCloser(body)
	}
	// The host's colon:port should be normalized. See Issue 14836.
	u.Host = removeEmptyPort(u.Host)
	req := &Request{
		ctx:        ctx,
		Method:     method,
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(Header),
		Body:       rc,
		Host:       u.Host,
	}
	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			req.ContentLength = int64(v.Len())
			buf := v.Bytes()
			req.GetBody = func() (io.ReadCloser, error) {
				r := bytes.NewReader(buf)
				return io.NopCloser(r), nil
			}
		case *bytes.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return io.NopCloser(&r), nil
			}
		case *strings.Reader:
			req.ContentLength = int64(v.Len())
			snapshot := *v
			req.GetBody = func() (io.ReadCloser, error) {
				r := snapshot
				return io.NopCloser(&r), nil
			}
		default:
		if req.GetBody != nil && req.ContentLength == 0 {
			req.Body = NoBody
			req.GetBody = func() (io.ReadCloser, error) { return NoBody, nil }
		}
	}

	return req, nil
}
```

整个构建请求过程中只有处理body的时候稍有一些复杂，我们需要根据body的类型使用不同的方法将其包装成`io.ReadCloser`类型；



### 启动事务

构建HTTP请求后，接下来我们需要开启HTTP事务进行请求并且等待远程响应，我们以net/http.Client.Do()方法为例子，我们看一下它的调用链路：

- net/http.Client.Do()
- net/http.Client.do()
- net/http.Client.send()
- net/http.Send()
- net/http.Transport.RoundTrip()

RoundTrip()是RoundTripper类型中的一个的方法，net/http.Transport是其中的一个实现，在net/http/transport.go文件中我们可以找到这个方法：

```go
// roundTrip implements a RoundTripper over HTTP.
func (t *Transport) roundTrip(req *Request) (*Response, error) {
	t.nextProtoOnce.Do(t.onceSetNextProtoDefaults)
	ctx := req.Context()
	trace := httptrace.ContextClientTrace(ctx)
    // 省略前置检查部分
    .....
	for {
        // 用来检测ctx退出信号
		select {
		case <-ctx.Done():
			req.closeBody()
			return nil, ctx.Err()
		default:
		}

		// 获取连接，这块就是我们要看的重点，Go语言通过连接池对资源进行了复用；
		pconn, err := t.getConn(treq, cm)
		if err != nil {
			t.setReqCanceler(cancelKey, nil)
			req.closeBody()
			return nil, err
		}

		var resp *Response
		if pconn.alt != nil {
			// HTTP/2 path.
			t.setReqCanceler(cancelKey, nil) // not cancelable with CancelRequest
			resp, err = pconn.alt.RoundTrip(req)
		} else {
            // 开始处理响应
			resp, err = pconn.roundTrip(treq)
		}
		if err == nil {
			resp.Request = origReq
			return resp, nil
		}

		// Rewind the body if we're able to.
		req, err = rewindBody(req)
		if err != nil {
			return nil, err
		}
	}
}
```

代码一大堆，我们只要重点看两部分即可：

- `net/http.Transport.getConn()`获取连接
- `net/http.persistConn.roundTrip（）`处理写入HTTP请求并在`select`中等待响应的返回；

### 获取连接

```go
func (t *Transport) getConn(treq *transportRequest, cm connectMethod) (pc *persistConn, err error) {
	req := treq.Request
	trace := treq.trace
	ctx := req.Context()
	if trace != nil && trace.GetConn != nil {
		trace.GetConn(cm.addr())
	}

	w := &wantConn{
		cm:         cm,
		key:        cm.key(),
		ctx:        ctx,
		ready:      make(chan struct{}, 1),
		beforeDial: testHookPrePendingDial,
		afterDial:  testHookPostPendingDial,
	}
	defer func() {
		if err != nil {
			w.cancel(t, err)
		}
	}()

	// 在队列中有闲置连接，直接返回
	if delivered := t.queueForIdleConn(w); delivered {
		pc := w.pc
		return pc, nil
	}

	cancelc := make(chan error, 1)
	t.setReqCanceler(treq.cancelKey, func(err error) { cancelc <- err })

	// 放到队列中等待建立新的连接
	t.queueForDial(w)

	// 阻塞等待连接
	select {
	case <-w.ready:
		return w.pc, w.err
	case <-req.Cancel:
		return nil, errRequestCanceledConn
	case <-req.Context().Done():
		return nil, req.Context().Err()
	case err := <-cancelc:
		if err == errRequestCanceled {
			err = errRequestCanceledConn
		}
		return nil, err
	}
}
```

因为连接的建议会消耗比较多的时间，带来较大的开下，所以Go语言使用了连接池对资源进行分配和复用，先调用net/http.Transport.queueForIdleConn()获取等待闲置的连接，如果没有获取到在调用`net/http.Transport.queueForDial`在队列中等待建立新的连接，通过select监听连接是否建立完毕，超时未获取到连接会上剖错误，我们继续在`queueForDial`追踪TCP连接的建立：

```go
func (t *Transport) queueForDial(w *wantConn) {
	w.beforeDial()
	if t.MaxConnsPerHost <= 0 {
		go t.dialConnFor(w)
		return
	}

	t.connsPerHostMu.Lock()
	defer t.connsPerHostMu.Unlock()

	if n := t.connsPerHost[w.key]; n < t.MaxConnsPerHost {
		if t.connsPerHost == nil {
			t.connsPerHost = make(map[connectMethodKey]int)
		}
		t.connsPerHost[w.key] = n + 1
		go t.dialConnFor(w)
		return
	}
    ....
}

```

我们会启动一个goroutine做tcp的建连，最终调用`dialConn`方法，在这个方法内做持久化连接，调用`net`库的dial方法进行TCP连接：

```go
func (t *Transport) dialConn(ctx context.Context, cm connectMethod) (pconn *persistConn, err error) {
	pconn = &persistConn{
		t:             t,
		cacheKey:      cm.key(),
		reqch:         make(chan requestAndChan, 1),
		writech:       make(chan writeRequest, 1),
		closech:       make(chan struct{}),
		writeErrCh:    make(chan error, 1),
		writeLoopDone: make(chan struct{}),
	}

	conn, err := t.dial(ctx, "tcp", cm.addr())
	if err != nil {
		return nil, err
	}
	pconn.conn = conn

	pconn.br = bufio.NewReaderSize(pconn, t.readBufferSize())
	pconn.bw = bufio.NewWriterSize(persistConnWriter{pconn}, t.writeBufferSize())

	go pconn.readLoop()
	go pconn.writeLoop()
	return pconn, nil
}
```

在连接建立后，代码中我们我们还看到分别启动了两个goroutine，`readLoop`用于从tcp连接中读取数据，`writeLoop`用于从tcp连接中写入数据；

我们看一下writeLoop方法：

```go
func (pc *persistConn) writeLoop() {
	defer close(pc.writeLoopDone)
	for {
		select {
		case wr := <-pc.writech:
			startBytesWritten := pc.nwrite
			err := wr.req.Request.write(pc.bw, pc.isProxy, wr.req.extra, pc.waitForContinue(wr.continueCh))
			if bre, ok := err.(requestBodyReadError); ok {
				err = bre.error
				wr.req.setError(err)
			}
			if err == nil {
				err = pc.bw.Flush()
			}
			if err != nil {
				if pc.nwrite == startBytesWritten {
					err = nothingWrittenError{err}
				}
			}
			pc.writeErrCh <- err // to the body reader, which might recycle us
			wr.ch <- err         // to the roundTrip function
			if err != nil {
				pc.close(err)
				return
			}
		case <-pc.closech:
			return
		}
	}
}
```

监听writech通道，所以的数据发送都是在这个循环中写入的；

`net/http.Transport{}`中提供了连接池配置参数，开发者可以自行定义：

```go
type Transport struct {
	......
    MaxIdleConns int
    MaxIdleConnsPerHost int
    MaxConnsPerHost int
    IdleConnTimeout time.Duration
    ......
}
```



### 处理HTTP请求

`net/http.persistConn.roundTrip()`会处理HTTP请求，我们看其具体实现：

```go
func (pc *persistConn) roundTrip(req *transportRequest) (resp *Response, err error) {
    ......
	startBytesWritten := pc.nwrite
	writeErrCh := make(chan error, 1)
    //writeLoop循环写入
	pc.writech <- writeRequest{req, writeErrCh, continueCh}
    
    // 等待响应数据
    resc := make(chan responseAndError)
    // readLoop循环等待响应
	pc.reqch <- requestAndChan{
		req:        req.Request,
		cancelKey:  req.cancelKey,
		ch:         resc,
		addedGzip:  requestedGzip,
		continueCh: continueCh,
		callerGone: gone,
	}
    
	for {
		select {
		case err := <-writeErrCh:
		case <-pcClosed:
			pcClosed = nil
			if canceled || pc.t.replaceReqCanceler(req.cancelKey, nil) {
				if debugRoundTrip {
					req.logf("closech recv: %T %#v", pc.closed, pc.closed)
				}
				return nil, pc.mapRoundTripError(req, startBytesWritten, pc.closed)
			}
		case <-respHeaderTimer:
			if debugRoundTrip {
				req.logf("timeout waiting for response headers.")
			}
			pc.close(errTimeout)
			return nil, errTimeout
		case re := <-resc:
			if (re.res == nil) == (re.err == nil) {
				panic(fmt.Sprintf("internal error: exactly one of res or err should be set; nil=%v", re.res == nil))
			}
			if debugRoundTrip {
				req.logf("resc recv: %p, %T/%#v", re.res, re.err, re.err)
			}
			if re.err != nil {
				return nil, pc.mapRoundTripError(req, startBytesWritten, re.err)
			}
			return re.res, nil
		case <-cancelChan:
			canceled = pc.t.cancelRequest(req.cancelKey, errRequestCanceled)
			cancelChan = nil
		case <-ctxDoneChan:
			canceled = pc.t.cancelRequest(req.cancelKey, req.Context().Err())
			cancelChan = nil
			ctxDoneChan = nil
		}
	}
}
```

我们重点关注这两个通道：

- `pc.writech` ：其类型是`chan writeRequest` ，writeLoop协程会循环写入数据，`net/http.Request.write`会根据`net/http.Request`结构中的字段按照HTTP协议组成TCP数据段，TCP协议栈会负责将HTTP请求中的内容发送到目标服务器上；
- `pc.reqch`：其类型是`chan requestAndChan`，readLoop协程会循环读取响应数据并且调用`net/http.ReadResponse`进行协议解析，其中包含状态码、协议版本、请求头等内容；



### 小结

我们简单总结一下net/http库中HTTP客户端的实现：

- net/http.Client是级别最高的抽象，其中`transport`用于开启HTTP事务，`jar`用于处理cookie；
- net/http.Transport中主要逻辑两部分：
  - 从连接池中获取持久化连接
  - 使用持久化连接处理HTTP请求

net/http库中默认有一个DefaultClient可以直接使用，DefaultClient有对应DefaultTransport，可以满足我们大多数场景，如果需要使用自己管理HTTP客户端的头域、重定向等策略，那么可以自定义Client，如果需要管理代理、TLS配置、连接池、压缩等设置，可以自定义Transport；

因为HTTP协议的版本是不断变化的，所以为了可扩展性，transport是一个接口类型，具体的是实现是`Transport`、`http2Transport`、`fileTransport`，这样实现扩展性变高，值得我们学习；

HTTP在建立连接时会耗费大量的资源，需要开辟一个goroutine去创建TCP连接，连接建立后会在创建两个goroutine用于HTTP请求的写入和响应的解析，然后使用channel进行通信，所以要合理利用连接池，避免大量的TCP连接的建立可以优化性能；



## 服务端

我们可以用`net/http`库快速搭建HTTP服务，HTTP服务端主要包含两部分：

- 注册处理器：net/http.HandleFunc函数用于注册处理器
- 监听端口：`net/http.ListenAndServe`用于处理请求



### 注册处理器

直接调用`net/http.HandleFunc`可以注册路由和处理函数：

```go
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}
```

我们可以看到处理函数是一个统一的格式：

```go
handler func(ResponseWriter, *Request)
```

默认调用`HTTP`服务起的`DefaultServeMux`处理请求，`DefaultServeMux`本质是`ServeMux`：

```go
type ServeMux struct {
	mu    sync.RWMutex  // 读写锁，保证并发安全，注册处理器时会加写锁做保护
	m     map[string]muxEntry // 路由规则，一个string对应一个mux实体，这里的string就是注册的路由表达式
	es    []muxEntry // slice of entries sorted from longest to shortest.
	hosts bool       // whether any patterns contain hostnames
}
```

- `mu`：需要加读写锁保证并发安全，注册处理器时会加写锁保证写map的数据正确性，这个map就是pattern和handler；
- `m`：存储路由规则，key就是pattern，value是muEntry实体，muEntry实体中包含：pattern和handler
- `es`：存储的也是muxEntry实体，因为我们使用map存储路由和handler的对应关系，所以只能索引静态路由，并不支持[path_param]，所以这块的作用是当在map中没有找到匹配的路由时，会遍历这个切片进行前缀匹配，这个切片按照路由长度进行排序；
- `hosts`：这个也是用来应对特殊case，如果我们注册的路由没有以`/`开始，那么就认为我们注册的路由包含host，所以路由匹配时需要加上host；

我看看一下路由注册函数：

```go
func (mux *ServeMux) Handle(pattern string, handler Handler) {
    // 加锁，保证并发安全
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if _, exist := mux.m[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
	e := muxEntry{h: handler, pattern: pattern}
    // map存储路由和处理函数的映射
	mux.m[pattern] = e
    // 如果路由最后加了`/`放入到切片后在路由匹配时做前缀匹配
	if pattern[len(pattern)-1] == '/' {
		mux.es = appendSorted(mux.es, e)
	}
	// 如果路由第一位不是/,则认为注册的路由加上了host，所以在路由匹配时使用host+path进行匹配；
	if pattern[0] != '/' {
		mux.hosts = true
	}
}
```



### 监听端口

`net/http`库提供了`ListenAndServe()`用来监听TCP连接并处理请求：

```go
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}
```

在这里初始化Server结构，然后调用`ListenAndServe`：

```go
func (srv *Server) ListenAndServe() error {
	if srv.shuttingDown() {
		return ErrServerClosed
	}
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}
    // 调用net进行tcp连接
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}
```

我们调用net网络库进行tcp连接，这里包含了创建socket、bind绑定socket与地址，listen端口的操作，最后调用Serve方法循环等待客户端的请求：

```go
func (srv *Server) Serve(l net.Listener) error {

	origListener := l
	l = &onceCloseListener{Listener: l}
	defer l.Close()

	if err := srv.setupHTTP2_Serve(); err != nil {
		return err
	}

	if !srv.trackListener(&l, true) {
		return ErrServerClosed
	}
	defer srv.trackListener(&l, false)

	baseCtx := context.Background()
	if srv.BaseContext != nil {
		baseCtx = srv.BaseContext(origListener)
		if baseCtx == nil {
			panic("BaseContext returned a nil context")
		}
	}

	var tempDelay time.Duration // how long to sleep on accept failure

	ctx := context.WithValue(baseCtx, ServerContextKey, srv)
	for {
        // 接收客户端请求
		rw, err := l.Accept()
		if err != nil {
			select {
			case <-srv.getDoneChan():
				return ErrServerClosed
			default:
			}
            // 网络错误进行延时等待
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				srv.logf("http: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		connCtx := ctx
		if cc := srv.ConnContext; cc != nil {
			connCtx = cc(connCtx, rw)
			if connCtx == nil {
				panic("ConnContext returned nil")
			}
		}
		tempDelay = 0
        // 创建一个新的连接
		c := srv.newConn(rw)
		c.setState(c.rwc, StateNew, runHooks) // before Serve can return
        // 读取起一个goroutine处理客户端请求
		go c.serve(connCtx)
	}
}
```

从上述代码我们可以到每个HTTP请求服务端都会单独创建一个goroutine来处理请求，我们一下处理过程：

```go
// Serve a new connection.
func (c *conn) serve(ctx context.Context) {
	c.remoteAddr = c.rwc.RemoteAddr().String()
	ctx = context.WithValue(ctx, LocalAddrContextKey, c.rwc.LocalAddr())
	var inFlightResponse *response
	defer func() {
        // 添加recover函数防止panic引发主程序挂掉；
		if err := recover(); err != nil && err != ErrAbortHandler {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			c.server.logf("http: panic serving %v: %v\n%s", c.remoteAddr, err, buf)
		}
	}()


	// HTTP/1.x from here on.
	ctx, cancelCtx := context.WithCancel(ctx)
	c.cancelCtx = cancelCtx
	defer cancelCtx()

	c.r = &connReader{conn: c}
	c.bufr = newBufioReader(c.r)
	c.bufw = newBufioWriterSize(checkConnErrorWriter{c}, 4<<10)

	for {
        // 读取请求，从连接中获取HTTP请求并构建一个实现了`net/http.Conn.ResponseWriter`接口的变量`net/http.response`
		w, err := c.readRequest(ctx)
		if c.r.remain != c.server.initialReadLimitSize() {
			c.setState(c.rwc, StateActive, runHooks)
		}
		if err != nil {
		}
		// 处理请求
		serverHandler{c.server}.ServeHTTP(w, w.req)
	}
}
```

我们继续跟踪`ServeHTTP`方法，ServeMux是一个HTTP请求的多路复用器，在这里可以根据请求的URL匹配合适的处理器，我们看代码：

```go
func (mux *ServeMux) ServeHTTP(w ResponseWriter, r *Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(StatusBadRequest)
		return
	}
    // 进行路由匹配，获取注册的处理函数
	h, _ := mux.Handler(r)
    // 这块就是执行我们注册的handler，也就是例子中的getProfile()
	h.ServeHTTP(w, r)
}
```

在`mux.Handler()`中我们就看到了路由匹配的代码：

```go
func (mux *ServeMux) handler(host, path string) (h Handler, pattern string) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	// Host-specific pattern takes precedence over generic ones
	if mux.hosts {
		h, pattern = mux.match(host + path)
	}
	if h == nil {
		h, pattern = mux.match(path)
	}
	if h == nil {
		h, pattern = NotFoundHandler(), ""
	}
	return
}
func (mux *ServeMux) match(path string) (h Handler, pattern string) {
	// 先从map中查找
	v, ok := mux.m[path]
	if ok {
        // 找打了返回注册的函数
		return v.h, v.pattern
	}

	// 从切片中进行前缀匹配
	for _, e := range mux.es {
		if strings.HasPrefix(path, e.pattern) {
			return e.h, e.pattern
		}
	}
	return nil, ""
}
```



### 小结

服务端的代码看主逻辑主要是看两部分，一个是注册处理器，标准库使用map进行存储，本质是一个静态索引，同时维护了一个切片，用来做前缀匹配，只要以`/`结尾的，都会在切片中存储；服务端监听端口本质也是使用net网络库进行TCP连接，然后监听对应的TCP连接，每一个HTTP请求都会开一个goroutine去处理请求，所以如果有海量请求，会在一瞬间创建大量的goroutine，这个可能是一个性能瓶颈点，所以小伙伴要注意下这块的性能问题；



## 总结

net/HTTP的总体代码行数是比较多的，我们只需要看主要逻辑是怎么实现的就可以了，别人问你原理能打出来个所以然就行，不必要扣细节，当出现问题或者想具体了解某部分协议的时候在细看源码对应部分即可；

像我们过了一遍源码后我们就知道当前`net/http`的一些优缺点，比如优点是HTTP客户端使用了连接池，避免频繁建立带来的大开销，缺点是HTTP服务端的路由只是一个静态索引匹配，对于动态路由匹配支持的不好，并且每一个请求都会创建一个gouroutine进行处理，海量请求到来时需要考虑这块的性能瓶颈；

net/http标准库可以让我们很快的就实现一个HTTP服务器，并且也有很多我们值得借鉴学习的地方，所以源码的学习还是很必要的；

**好啦，今天的文章就到这里了，我是`asong`，我们下期见~；**

![](https://song-oss.oss-cn-beijing.aliyuncs.com/golang_dream/article/static/%E6%89%AB%E7%A0%81_%E6%90%9C%E7%B4%A2%E8%81%94%E5%90%88%E4%BC%A0%E6%92%AD%E6%A0%B7%E5%BC%8F-%E7%99%BD%E8%89%B2%E7%89%88.png)

