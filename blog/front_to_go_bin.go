package main

// https://my.oschina.net/henrylee2cn/blog/741372
// 教你如何将前端文件打包进Go程序

// 在Golang的开发中，我们有时会想要将一些外部依赖文件打包进二进制程序。比如本人在开发lessgo web框架时，希望将扩展包swagger（一个自动API文档的前端）打包进项目文件中，从而减少依赖，并能提高代码稳定性。实现步骤如下：
//
// 下载两个Golang的第三方包
// go get github.com/jteeuwen/go-bindata/...
// go get github.com/elazarl/go-bindata-assetfs/...
//
//
// 使用 “go install” 命令分别编译获得 go-bindata.exe 和 go-bindata-assetfs.exe 文件
//
// 执行 “go-bindata-assetfs.exe views/...” 将./views目录下所有文件写入 bindata_assetfs.go 文件
//
// bindata_assetfs.go文件中提供了名为 assetFS() 的函数，它返回包含了view文件内容的 http.Filesystem 接口实例
//
// 以静态文件路由为例，调用方式为：
//
// http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(assetFS)))

func main() {
	
}
