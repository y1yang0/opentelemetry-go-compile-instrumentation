# build instrumentation tool
go build -a -o otel
TOOL=$(pwd)/otel
# compile-time instrumentation via toolexec
cd demo
go build -a -toolexec=$TOOL -o demo .
./demo