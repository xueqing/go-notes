# go 开发遇到的坑

记录 go 开发遇到的坑

- [int 类型在 64 位机器是 8 字节](int.md)
- [json 解析成结构体](json_struct.md)
- [module 不要在 vscode 工作区打开工程](mod_workspace.md)
- [同一包不同源文件变量声明时不要带包名](package_var.md)
- [string 类型](string.md)
- [不使用短变量声明 `:=`](var_scope.md)
- [cgo 调用宏函数](cgo_macro_func.md)
- [defer 语句需要注意的](defer.md)
- [WaitGroup 不是引用类型](waitgroup.md)

- 有时候新定义的变量或者函数不能跳转，需要重启 vscode，相关插件长时间运行可能崩溃
