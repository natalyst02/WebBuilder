currentDir := $(shell pwd)

mock:
	cd $(currentDir)
	echo "Directory: $(currentDir)"
	rm -rf mocks
	mockery --all --dir ./media --outpkg mock_media --output testing/mocks/mock_media
	mockery --all --dir ./content --outpkg mock_content --output testing/mocks/mock_content