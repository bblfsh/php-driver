test-native-internal:
	composer.phar install
	echo "test-native-internal > not implemented"

build-native-internal:
	cp -rf native $(BUILD_PATH)/src
	cd $(BUILD_PATH)/src && composer.phar install

	rm $(BUILD_PATH)/bin/native || true
	ln -s /opt/driver/src/ast $(BUILD_PATH)/bin/native

include .sdk/Makefile
