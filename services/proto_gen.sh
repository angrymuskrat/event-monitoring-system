protoc --gofast_out=plugins=grpc:. proto/data.proto
protoc --gofast_out=plugins=grpc:. event-detection/proto/service.proto
protoc --gofast_out=plugins=grpc:. data-storage/proto/data-storage.proto

# service do not compile due to data-storage/proto/data-storage.pb.go cannot see "data", but with absolute go module path it worksbash
perl -pi -e 's/proto1 "proto"/proto1 "github.com\/angrymuskrat\/event-monitoring-system\/services\/proto"/g' data-storage/proto/data-storage.pb.go
