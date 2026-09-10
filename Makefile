libfio.so libfio.h:
	go build -buildmode=c-shared -ldflags="-s -w" -o libfio.so ./cbindings

clean:
	rm -f libfio.so libfio.h

