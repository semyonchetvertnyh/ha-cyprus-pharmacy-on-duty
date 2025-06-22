BINARY=homeassistant/packages/pharmacy/on-duty/bin
TARGET=root@homeassistant.local:/homeassistant/packages/pharmacy/on-duty

build:
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -o $(BINARY) ./app

deploy: build
	scp ./$(BINARY) $(TARGET)

clean:
	rm -f $(BINARY)