# 安装 go

- [安装 go](#安装-go)
  - [Linux 安装和使用 go](#linux-安装和使用-go)
    - [安装](#安装)
    - [设置工作目录 GOPATH](#设置工作目录-gopath)
    - [测试安装](#测试安装)
    - [安装其他版本](#安装其他版本)
  - [Windows 安装和使用 go](#windows-安装和使用-go)
  - [MacOS 安装和使用 go](#macos-安装和使用-go)
  - [卸载旧版本](#卸载旧版本)
  - [命令](#命令)
  - [vscode 使用 go](#vscode-使用-go)
    - [配置代理](#配置代理)
    - [安装插件](#安装插件)
    - [私有仓库使用](#私有仓库使用)
    - [使用 cgo](#使用-cgo)
  - [参考](#参考)

## Linux 安装和使用 go

### 安装

- 下载[安装包](https://golang.org/dl/)
  - 选择最新的 [Linux 版本](https://dl.google.com/go/go1.11.2.linux-amd64.tar.gz)
  - 下载`wget https://dl.google.com/go/go1.11.2.linux-amd64.tar.gz`
- [安装](https://golang.org/doc/install)
  - 解压：`tar -C /usr/local -zxf go1.11.2.linux-amd64.tar.gz`
- 配置环境变量
  - 修改`/etc/profile`或`~/.profile`
    - 追加`export PATH=$PATH:/usr/local/go/bin`
    - 执行`source`命令更新配置文件立即生效

### 设置工作目录 GOPATH

- 默认工作目录是 `$HOME/go`
- 工作目录下面有三个文件夹
  - src：存放源码的目录，新建项目都在该目录下
  - pkg：编译生成的包文件存放目录
  - bin：编译生成的可执行文件和 go 相关的工具
- 如果需要自定义工作目录：
  - **建议：**不要和 go 的安装目录相同
  - 修改 `~/.bashrc`，添加 `export GOPATH=$HOME/go`
  - 执行 `source ~/.bashrc` 使脚本生效
  - 修改`/etc/profile`或`~/.profile`，添加 `export GOROOT=$HOME/go`，将 `$HOME/go/bin` 加入系统环境变量 `PATH`
  - 执行`source`命令更新配置文件立即生效
  
### 测试安装

- 创建并进入默认工作目录 `~/go`
- 创建并进入目录 `src/hello`
- 创建文件 `hello.go`

  ```go
  package main

  import "fmt"

  func main() {
    fmt.Printf("Hello, world\n")
  }
  ```

- 编译：`go build`，生成可执行文件 `hello`
- 运行：`./hello`
- 安装二进制文件到工作目录的 `bin` 目录：`go install`
- 删除工作目录的二进制文件：`go clean -i`

### 安装其他版本

- 安装版本 1.10.7
  - `go get golang.org/dl/go1.10.7`
  - `go1.10.7 download`
- 使用版本 1.10.7
  - `go1.10.7 version`

## Windows 安装和使用 go

- 下载 Windows [安装包](https://golang.org/dl/) msi 文件
- 安装到 `c:\Go` 目录
- 将 `c:\Go\bin` 加入系统环境变量 `PATH`
- 测试安装
  - 创建 go 的工作目录，比如 `g:\gopro`
  - 设置工作目录路径：在用户变量中加入 `GOPATH`
  - 将 `g:\gopro\bin` 加入系统环境变量 `PATH`
  - 创建 `g:\gopro\src\hello` 目录，创建 `hello.go` 文件
  - 打开 Windows 终端，切换到 `g:\gopro\src\hello` 目录
    - 编译：`go build`，生成可执行文件 `hello.exe`
    - 运行：`hello`
    - 安装二进制文件到工作目录的 `bin` 目录：`go install`
    - 删除工作目录的二进制文件：`go clean -i`

## MacOS 安装和使用 go

- 执行 `brew install go`
- 配置环境变量
  - 修改 `~/.profile`
    - 追加`export GOROOT=/usr/local/go`
    - 追加`export PATH=$PATH:$GOROOT/bin`
    - 执行`source ~/.profile`命令更新配置文件立即生效
- 验证配置：执行 `go version`
- 写测试程序
  - 写 `hello.go`

    ```go
    package main

    import "fmt"

    func main() {
      fmt.println("Hello World!")
    }
    ```

  - 执行 `go run hello.go`

## 卸载旧版本

- 删除 go 目录
  - Linux/MacOS/FreeBSD `/usr/local/go`
  - Windows `c:\go`
- 从环境变量 `PATH` 中删除 go 的 bin 目录
  - Linux/FreeBSD 编辑 `/etc/profile` 或 `~/.profile`
  - MacOS 删除 `/etc/paths.d/go`

## 命令

- 查看 golang 环境变量 `go env`
- 公开模块导入超时，配置代理 `go env -w GOPROXY="https://goproxy.io"`

## vscode 使用 go

### 配置代理

- 配置环境变量 `go env -w GOPROXY="https://goproxy.io"`
- 重启 vscode

### 安装插件

- 安装插件 `Go`
- `Ctrl+Shift+P`，输入 `go`
- 选择 `Install/Update Tools`
- 全选，安装
- 重启 vscode

### 私有仓库使用

如果有一个私有仓库 <gitlab.bmi>，使用 `go get` 获取仓库时会报错：

```sh
go get -v gitlab.bmi/ylrc/bmi-av/rtspproxy.git
go get gitlab.bmi/ylrc/bmi-av/rtspproxy.git: module gitlab.bmi/ylrc/bmi-av/rtspproxy.git: reading https://goproxy.io/gitlab.bmi/ylrc/bmi-av/rtspproxy.git/@v/list: 404 Not Found
  server response:
  not found: go list -m -json -versions gitlab.bmi/ylrc/bmi-av/rtspproxy.git@latest
  :  
```

`go get` 通过带来无法访问私有仓库，因此出现 404 错误。需要配置 `GOPRIVATE` 环境变量，指定域名为私有仓库。

```sh
go env -w GOPRIVATE="*.bmi"
```

再次拉取源码仍然报错

```sh
go get -v gitlab.bmi/ylrc/bmi-av/rtspproxy.git
# cd .; git ls-remote https://gitlab.bmi/ylrc/bmi-av/rtspproxy
fatal: unable to access 'https://gitlab.bmi/ylrc/bmi-av/rtspproxy/': Failed to connect to gitlab.bmi port 443: Connection refused
kiki@gitlab.bmi's password:
```

这是因为默认配置 ssh 公钥访问私有仓库，因此需要配置 git 拉取私有仓库时使用 ssh 而不是 https。

```sh
# 方法 1：使用 git 命令配置
git config --global url."git@gitlab.bmi".insteadOf https://gitlab.bmi
## 或者
git config --global url."ssh://git@gitlab.bmi".insteadOf https://gitlab.bmi
# 方法 2：修改 git 配置文件
vim ~/.gitconfig
## 加入下面的内容
[url "git@gitlab.bmi:"]
  insteadOf = https://gitlab.bmi
## 或
[url "ssh://git@gitlab.bmi/"]
  insteadOf = https://gitlab.bmi/
```

### 使用 cgo

比如使用 ffmpeg 库进行开发

```sh
# 编译 ffmpeg：安装目录是 /home/kiki/ffmpeg/ffmpeg-4.1
# 添加 ffmpeg 动态库路径
vim ~/.bashrc
# export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/home/kiki/ffmpeg/ffmpeg-4.1/lib
source ~/.bashrc
# 更新缓存
sudo ldconfig
# 修改 go 环境变量
go env -w CGO_CFLAGS="-I/home/kiki/ffmpeg/ffmpeg-4.1/include"
go env -w CGO_LDFLAGS="-L/home/kiki/ffmpeg/ffmpeg-4.1/lib -lavcodec -lavformat -lavutil -lswscale -lswresample -lavdevice -lavfilter"
```

## 参考

- [Go填坑之将Private仓库用作module依赖](https://segmentfault.com/a/1190000021127791)
