docker-compose up -d

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

验证 redis 连接
docker exec -it redis-cms redis-cli ping

redis-cli
get session_id:admin
