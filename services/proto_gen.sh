protoc --gofast_out=plugins=grpc:. proto/data.proto
protoc --gofast_out=plugins=grpc:. dbsvc/proto/dbsvc.proto

# service don't compilate due to dbsvc/proto/dbsvc.pb.go cann't see "data", but with absolute go module path it works
# maybe I have incorrect enviroment
perl -pi -e 's/proto1 "proto"/proto1 "github.com\/angrymuskrat\/event-monitoring-system\/services\/proto"/g' dbsvc/proto/dbsvc.pb.go
