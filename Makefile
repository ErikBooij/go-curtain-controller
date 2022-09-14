build-local:
	go build -o dist/curtain-controller-macos main.go && chmod +x dist/curtain-controller-macos

build-remote:
	GOOS=linux GOARCH=amd64 go build -o dist/curtain-controller *.go && chmod +x dist/curtain-controller

deploy:
	@echo "Building application"
	@make build-remote > /dev/null
	@echo "Stopping service"
	@make remote-stop > /dev/null
	@echo "Backing up existing binary"
	@ssh server cp /home/erikbooij/curtain-controller curtain-controller.bak || true
	@echo "Uploading binary"
	@scp dist/curtain-controller server:/home/erikbooij/curtain-controller > /dev/null
	@echo "Starting service"
	@make remote-start > /dev/null
	@echo "Waiting for application to boot"
	@sleep 3 > /dev/null
	@echo "Verify service is running"
	@make remote-verify-running > /dev/null
	@echo "Verify expected output"
	@make remote-verify-success > /dev/null

remote-start:
	ssh server sudo systemctl restart curtain-controller || true

remote-stop:
	ssh server sudo systemctl stop curtain-controller || true

remote-verify-running:
	ssh server sudo systemctl is-active curtain-controller --quiet

remote-verify-success:
	ssh server sudo journalctl -u curtain-controller --since \"12 seconds ago\" | grep "running"

run-local: build-local
	./dist/curtain-controller-macos