# 基础镜像
FROM ubuntu:20.04
# 把编译后的可执行文件复制到工作目录
COPY webook /app/webook
# 切换到工作目录
WORKDIR /app
# 执行命令(也可以使用CMD)
ENTRYPOINT ["/app/webook"]