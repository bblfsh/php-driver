build-native-internal:
	echo "build-native-internal > not implemented"

test-native-internal:
	echo "test-native-internal > not implemented"

test-native:
	cd native; \
	echo "test-native > not implemented"

build-native:
	cd native; \
	echo "build-native > not implemented"
	echo -e "#!/bin/bash\necho 'not implemented'" > $(BUILD_PATH)/native
	chmod +x $(BUILD_PATH)/native

include .sdk/Makefile