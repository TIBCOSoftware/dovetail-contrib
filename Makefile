SCRIPTS_PATH      := scripts

.PHONY: all
all: functions

.PHONY: functions
functions: 
	$(SCRIPTS_PATH)/functions.sh ${IMAGE_NAME} ${IMAGE_TAG} ${IMAGE_URL}