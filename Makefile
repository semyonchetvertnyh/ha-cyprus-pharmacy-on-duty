BINARY=./homeassistant/packages/pharmacy/on-duty/bin
TARGET=/config/packages/pharmacy/on-duty/bin
TARGET_HOST=root@homeassistant.local

build:
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -o $(BINARY) ./app
	chmod +x $(BINARY)

deploy: build
	cat $(BINARY) | ssh $(TARGET_HOST) "cat > $(TARGET)"

clean:
	rm -f $(BINARY)