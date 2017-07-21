go build -ldflags " \
-X main.SERVER_VERSION=1.0.0 \
-X 'main.BUILD_TIME=`date`' \
-X 'main.GO_VERSION=`go version`' \
-X 'main.BEEGO_VERSION=`bee version | grep beego`' \
-X 'main.GIT_ADDR=`git remote -v | grep fetch | awk '{print $2}'`' \
-X 'main.GIT_COMMIT=`git rev-parse HEAD`'" 
