echo '删除日志文件'
rm -f nsqlookupd.log
rm -f nsqd.dat nsqd.log
rm -f nsqadmin.log
rm -f log1.dat log2.dat log3.dat
rm -f nohup.out

echo "启动MongoDB"
nohup mongod 2>&1&

echo '启动nsqlookupd服务'
nohup nsqlookupd >nsqlookupd.log 2>&1&

echo '启动nsqd服务'
nohup nsqd --lookupd-tcp-address=0.0.0.0:4160 -tcp-address="0.0.0.0:4150" > nsqd.log 2>&1&

echo "启动chat server"
go build
nohup ./nsq-chat > log1.dat 2>&1&

echo "启动archive"
nohup go run archive/archive.go > log2.dat 2>&1&

echo "启动bot"
nohup go run bot/bot.go > log3.dat 2>&1&
