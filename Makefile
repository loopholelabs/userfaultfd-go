all: minio

minio:
	podman run --name 'userfaultfd-go-minio' -d -p '9000:9000' -p '9001:9001' 'minio/minio' server /data --console-address ':9001'

clean:
	podman rm -f 'userfaultfd-go-minio'
