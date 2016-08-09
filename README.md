mystories service模块
=======

mystories是一个新闻阅读类app服务端服务层的demo,这是用go语言实现的一个rpc服务。
rpc框架拷贝自go语言内置的net/rpc, 代码有改动。rpc的编码协议使用的是yar/msgpack。

安装
------------

可以直接通过go来下载安装代码:

```
go get -u github.com/mystories/service
```

或者直接下载源码:

```
git clone https://github.com/mystories/service.git
```

运行
------------

```
export GOPATH=~/work/go:~/work/mystories/service
go build src/main/main.go
```

其中GOPATH及源码路径需要调整为你自己的路径,参考build.sh脚本。

文档
------------

* Article.List:获取文章列表

* Article.Get:获取文章详情
