# seedHabits
种子习惯
这是基于go语言开发的一个项目。使用gin+xorm+postgres框架。

使用了swag自动生成文档，可在http://localhost:8001/swagger/index.html页面上查看及测试。
swag相关使用可查看。

使用jager做链路追踪，可在http://localhost:16686/search查看接口信息详情。
jagger相关使用可查看[Jaeger是干啥的](https://github.com/aogoZY/CodeExerciseDemo/blob/master/golang/jaeger%E5%88%9D%E8%AF%86.md)

使用了pprof做性能分析，可在http://localhost:8080/ui/source查看具体关注项指标。
pprof相关使用可查看[pprof--你想要的数据分析我都有](https://github.com/aogoZY/CodeExerciseDemo/blob/master/golang/pprof%E4%BD%BF%E7%94%A8%E8%AF%A6%E8%A7%A3-%E6%80%A7%E8%83%BD%E5%88%86%E6%9E%90%E5%A4%A7%E6%9D%80%E5%99%A8.md)

sdk里面集成了一些常用的三方库，比如加解密包、log日志包、redis使用demo、elasticsearch demo等。

这个项目是仿照以前我用的一个app种子习惯写的，主要有三大模块。

用户模块、记账模块、打卡模块。

支持用户创建自己的习惯，每日打卡。帮助每个人构建自己的时光机和回忆。

正在努力完善ing..
