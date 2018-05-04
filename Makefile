-include .sdk/Makefile

$(if $(filter true,$(sdkloaded)),,$(error You must install bblfsh-sdk))

test-native-internal:
	cd native && \
	composer.phar install && \
	./vendor/bin/phpunit tests/

build-native-internal:
	cp -rf native $(BUILD_PATH)/src
	cd $(BUILD_PATH)/src && composer.phar install

	ls -la $(BUILD_PATH)/bin/native
	rm $(BUILD_PATH)/bin/native || true
	ln -s /opt/driver/src/ast $(BUILD_PATH)/bin/native
	ls -la $(BUILD_PATH)/bin/native

