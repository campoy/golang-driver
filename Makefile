-include .sdk/Makefile

$(if $(filter true,$(sdkloaded)),,$(error You must install bblfsh-sdk))

test-native-internal:
	go test ./native

build-native-internal:
	go build -o $(BUILD_PATH)/bin/native ./native 