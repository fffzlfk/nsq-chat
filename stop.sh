echo "关闭MongoDB"
mongod --shutdown

echo "关闭nsq"
ps -ef | grep nsq | grep -v grep | awk '{print $2}'| xargs kill -2
ps -ef | grep nsqadmin | grep -v grep | awk '{print $2}'| xargs kill -2

echo "关闭服务"
ps -ef | grep ./nsq-chat | grep -v grep | awk '{print $2}'| xargs kill -9 -2
ps -ef | grep archive | grep -v grep | awk '{print $2}'| xargs kill -2
ps -ef | grep bot | grep -v grep | awk '{print $2}'| xargs kill -2

echo '删除日志文件'
rm -f *.log
rm -f *.dat
rm -f nohup.out

