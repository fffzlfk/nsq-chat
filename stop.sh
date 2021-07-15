echo "关闭MongoDB"
mongod --shutdown

echo "关闭nsq"
ps -ef | grep nsq | grep -v grep | awk '{print $2}'| xargs kill -2
ps -ef | grep nsqadmin | grep -v grep | awk '{print $2}'| xargs kill -2

echo "关闭服务"
ps -ef | nsq-chat | grep -v grep | awk '{print $2}'| xargs kill -2
ps -ef | grep archive | grep -v grep | awk '{print $2}'| xargs kill -2
ps -ef | grep bot | grep -v grep | awk '{print $2}'| xargs kill -2

echo '删除日志文件'
rm -f nsqlookupd.log
rm -f nsqd.dat nsqd.log
rm -f nsqadmin.log
rm -f log1.dat log2.dat log3.dat
rm -f nohup.out

