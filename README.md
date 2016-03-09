一个将你的应用运行为daemon的小工具

## 使用方法
要让你的应用支持daemon很简单，只需导入godaemon包即可，无需再调用任何方法，该包不导出任何方法和变量。

注意： 因为使用了syscall包，Win下不支持哦



```
package main
import(
    _ "github.com/tim1020/godaemon"
)
func main(){
	//正常的业务代码
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pid=%d", os.Getpid())
	})
	http.ListenAndServe(":8080", nil)
}
```
## 命令行参数

加入daemon支持后，你可以在命令行中以下列指令来管理你的应用:

（假设你编译出来的应用执行文件叫“app”,如果你的应用本身需要带参数运行，请把daemon指令参数放在最后）

- app start

	带start参数，启动为daemon进程

- app restart

	带start参数，重启进程，重启方式为发送一个kill -HUP 信号给旧的进程，再启动新的进程

	如果不存在正在运行的进程，则与start效果一样，直接启动一个daemon进程

- app stop

	带stop参数，停止应用进程

- app -h

	带-h参数， 显示命令行指令提示

- app

	不带参数或带其它不识别的参数启动，daemon不接管，直接短路至业务代码