SCRIPTS_PATH      := scripts

.PHONY: all
all: functions

.PHONY: release
release: hyperledger-fabric

.PHONY: tag
tag: 
	$(SCRIPTS_PATH)/tag.sh ${BUILD_BRANCH} ${BUILD_NUM}

.PHONY: is-prerelease
is-prerelease: 
	@ $(SCRIPTS_PATH)/prerelease.sh ${BUILD_BRANCH}

.PHONY: release-notes
release-notes: 
	@ $(SCRIPTS_PATH)/release-notes.sh

.PHONY: functions
functions: 
	$(SCRIPTS_PATH)/functions.sh ${IMAGE_NAME} ${IMAGE_TAG} ${IMAGE_URL}

.PHONY: hyperledger-fabric
hyperledger-fabric: 
	$(SCRIPTS_PATH)/hyperledger-fabric.sh