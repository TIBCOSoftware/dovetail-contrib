SCRIPTS_PATH      := scripts

.PHONY: all
all: functions

.PHONY: release
release: hyperledger-fabric

.PHONY: tag
tag: 
	$(SCRIPTS_PATH)/tag.sh

.PHONY: release_notes
release_notes: 
	$(SCRIPTS_PATH)/release_notes.sh

.PHONY: functions
functions: 
	$(SCRIPTS_PATH)/functions.sh ${IMAGE_NAME} ${IMAGE_TAG} ${IMAGE_URL}

.PHONY: hyperledger-fabric
hyperledger-fabric: 
	$(SCRIPTS_PATH)/hyperledger-fabric.sh