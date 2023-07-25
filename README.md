#  抖音简易后端
### 字节跳动青训营 GOGOG 队项目

## 前期准备进度
by 关竣佑 谢声儒

* 集成 GORM 框架
* 准备实体类 Model  
  所有的实体类均继承 CommonEntity (Id , CreateTime , IsDelete)
* 数据库设计  
  所有删除均采用逻辑删除 （规定删除后is_delete字段为 1 ）  
  所有的主键 id  采用bingint 对应的 go 实体类使用 int64 , 生成的时候采用分布式ID生成方式
  

