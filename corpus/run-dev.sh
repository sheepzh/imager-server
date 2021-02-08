segFileDirPath=/home/zhy/workspace/poetry/data/seged
logLevel=debug

go clean 

go build

if [ -f "./poem" ]; then
./poem -s=$segFileDirPath --log-level=$logLevel
fi
