docker-compose -p imooc-content-system up -d

# MySQL

进入 MySQL 容器
docker exec -it mysql-cms mysql -uroot -p
然后输入 root123 登录 MySQL。
show databases
use cms_account
show create table account

检查数据卷
docker volume ls
可以看到 mysql-data 卷已经创建。local imooc-content-system_mysql-data

MySQL 容器只在首次启动时（即数据卷为空时）执行/docker-entrypoint-initdb.d 目录中的脚本
如果容器已经启动过并且数据卷保留，再次启动容器不会重新执行初始化脚本
我们通过 docker-compose down -v 删除卷，确保再次启动时执行初始化脚本

重启 Docker 容器
docker-compose down -v
docker-compose up -d

验证数据库创建
docker exec mysql-cms mysql -uroot -proot123 -e "SHOW DATABASES;"

验证表创建
docker exec mysql-cms mysql -uroot -proot123 -e "USE cms_account; SHOW TABLES;"
docker exec mysql-cms mysql -uroot -proot123 -e "USE cms_content; SHOW TABLES;"

检查一下 SQL 脚本的挂载情况
docker exec mysql-cms ls -la /docker-entrypoint-initdb.d/

检查一下 MySQL 的启动日志
docker logs mysql-cms | grep -i "user.sql"

手动执行 SQL 脚本
docker exec -i mysql-cms mysql -uroot -proot123 < script/user.sql
docker exec -i mysql-cms mysql -uroot -proot123 < script/content.sql

# Redis

验证 redis 连接
docker exec -it redis-cms redis-cli ping

redis-cli
get session_id:admin

# DDD

在这个项目中，分层架构如下：

- Model 层：定义了与数据库表对应的结构体。负责：数据结构定义、表映射关系、简单的结构体方法
- DAO 层：封装所有数据库操作。负责：CRUD 操作、查询构建、事务管理
- Service 层：处理业务逻辑和控制流。负责：请求处理、业务逻辑、调用 DAO 层、返回响应

这样的分层设计有以下优点和缺点：
优点：
✅ 关注点分离：每一层有明确职责，便于维护
✅ 代码复用：DAO 层方法可被多个 Service 复用
✅ 测试便利：各层可独立测试，尤其是 DAO 层
✅ 变更隔离：数据库结构变更主要影响 Model 和 DAO 层
缺点/改进空间：
❌ Service 与 Controller 混合：当前设计中 Service 直接处理 HTTP 请求，职责不够清晰
❌ 缺少领域层抽象：业务规则分散在 Service 中，缺少领域模型抽象
❌ 事务处理不明确：复杂事务如何处理未在当前代码中体现

一般最佳实践是什么？
更理想的分层架构通常包括：

1. 清晰的职责分离
   Controller -> Service -> Repository -> Model
   ↓ ↓ ↓ ↓
   处理 HTTP 业务逻辑 数据访问 数据结构

2. 领域驱动设计(DDD)风格
   对于大型复杂系统，可考虑引入领域层：

Controller -> Application Service -> Domain Service -> Repository -> Model
↑
Domain Model

3. 具体改进建议

   - 分离 Controller 和 Service
   - Repository 模式替代 DAO
     更加面向领域，抽象数据访问
     可提供事务支持和更高级抽象
   - 使用 DTO 模式
     Model: 数据库模型
     DTO: 数据传输对象，用于 API
     Domain Model: 领域模型(可选)
   - 依赖注入
     使用依赖注入管理各层依赖关系
     便于测试和模块替换

4. 正确的分层方式

   - Controller 层应该包含：
     请求解析和验证
     调用 Service 层方法
     错误处理和 HTTP 响应构建
   - Service 层应该包含：
     业务逻辑处理
     调用 Repository/DAO 层
     事务管理
     领域规则验证

5. Key Differences Between DAO and Repository

   Repository Pattern: How It Would Look
   The Repository pattern is a higher-level abstraction focused on domain objects rather than database tables. Here's how it would differ:

   - Abstraction Level:
     DAO: Tightly coupled to database tables and operations (CRUD)
     Repository: Abstracts persistent storage, focuses on domain objects
   - Interface vs Implementation:
     DAO: Usually a concrete class with database operations
     Repository: Defines an interface first, implementations can vary
   - Object Mapping:
     DAO: Works directly with DB models
     Repository: Translates between domain objects and DB models
   - Domain Focus:
     DAO: Database-centric
     Repository: Domain-centric, part of the domain layer
   - Testability:
     Repository pattern is easier to mock due to its interface-based design

6. Benefits of Repository Over DAO
   - Better separation of concerns - domain logic doesn't know about database details
   - More flexible - can switch implementations (e.g., from MySQL to MongoDB) without changing the service layer
   - Easier testing - interfaces facilitate mocking
   - More domain-focused - represents collections of domain objects rather than database tables

The Repository pattern helps move from a data-centric application to a more domain-driven design, which is better suited for complex business logic.
