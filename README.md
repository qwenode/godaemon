为应用增加daemon和graceful


## 使用方法

- 只增加daemon功能

要让你的应用支持daemon很简单，只需导入godaemon包即可，无需再调用任何方法

```
package main
import(
    _ "github.com/tim1020/godaemon" //仅导入，包的init方法被自动调用，嵌入daemon功能
)
func main(){
	//正常的业务代码
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pid=%d", os.Getpid())
	})
	http.ListenAndServe(":8080", nil)
}
```

- 增加graceful

```
package main
import(
    "github.com/tim1020/godaemon" //注意： 与仅daemon时不一样
)
func main(){
	mux1 := http.NewServeMux()
	mux1.HandleFunc("/", handler)
    godaemon.GracefulServe(":8080", mux1)
}

func handler(w http.ResponseWriter, r *http.Request){
	//业务处理
}
````

> 注意： 因为使用了syscall包，Win下不支持哦


## 命令行操作

使用godaemon后，你可以在命令行中以下列指令来管理你的应用:

（假设你编译出来的应用执行文件叫“app”,如果你的应用本身需要带参数运行，请把daemon指令参数放在最后）

- app [start]

	带start参数或不带参数，启动为daemon进程

- app restart

	带restart参数，重启进程（GracefulServe时不中断服务，仅daemon时会先停止再启动）

- app stop

	带stop参数，停止应用进程 （kill -HUP）

- app -h

	带-h参数， 显示命令行指令提示


> 带其它不识别的参数启动，godaemon不接管，直接短路至业务代码
