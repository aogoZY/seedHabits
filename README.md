# seedHabits
种子习惯
这是基于go语言开发的一个项目。使用gin+xorm+postgres框架。

使用了swag自动生成文档，可在http://localhost:8001/swagger/index.html页面上查看及测试。

使用jager做链路追踪，可在http://localhost:16686/search查看接口信息详情。

sdk里面集成了一些常用的三方库，比如加解密包、log日志包、redis使用demo、elasticsearch demo等。

这个项目是仿照以前我用的一个app种子习惯写的，主要有三大模块。

用户模块、记账模块、打卡模块。

支持用户创建自己的习惯，每日打卡。帮助每个人构建自己的时光机和回忆。

正在努力完善ing。

