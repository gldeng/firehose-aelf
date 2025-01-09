rm -rf firehose-data Logs
mkdir firehose-data
mkdir Logs

docker run --platform linux/amd64 -it \
-p 10015:10015 \
-v ./start-aelf-node:/app/AElf.Launcher \
-v ./standard.yaml:/app/standard.yaml \
-v ./firehose-data:/app/firehose-data \
-v ./Logs:/app/Logs \
fireaelf \
-c /app/standard.yaml start


# --entrypoint "/app/fireaelf -c /app/standard.yaml" \

#docker run --platform linux/amd64 -it \
#-p 10015:10015 \
#--entrypoint /bin/bash \
#-v ./standard.yaml:/app/standard.yaml \
#-v ./firehose-data:/app/firehose-data \
#-v ./Logs:/app/Logs \
#fireaelf
#docker run --platform linux/amd64 -it \
#-p 10015:10015 \
#--entrypoint /bin/bash \
#-v ./start-aelf-node:/app/AElf.Launcher \
#-v ./standard.yaml:/app/standard.yaml \
#-v ./firehose-data:/app/firehose-data \
#-v ./Logs:/app/Logs \
#fireaelf