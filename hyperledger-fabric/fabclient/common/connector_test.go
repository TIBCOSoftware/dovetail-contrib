package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var connectionObject = `
{
	"id": "89905300-555c-11e9-98c5-0fe3500cf4b8",
	"type": "flogo:connector",
	"version": "1.0.0",
	"name": "fabclient-connector",
	"inputMappings": {},
	"outputMappings": {},
	"iteratorMappings": {},
	"title": "Fabric Connector",
	"description": "Fabric Connection",
	"ref": "/fabclient/connector/fabconnector",
	"settings": [
	 {
	  "name": "name",
	  "description": "Unique name of the Fabric network connection",
	  "type": "string",
	  "required": true,
	  "display": {
	   "name": "Name",
	   "visible": true,
	   "readonly": false,
	   "valid": true
	  },
	  "value": "local-first-network"
	 },
	 {
	  "name": "description",
	  "description": "Describe the Fabric network connection",
	  "type": "string",
	  "required": false,
	  "display": {
	   "name": "Description",
	   "visible": true,
	   "readonly": false,
	   "valid": true
	  },
	  "value": "Connection to Fabric sample first-network on localhost"
	 },
	 {
		"name": "entityMatcher",
		"type": "string",
		"required": false,
		"display": {
		 "name": "Connection entity matcher file",
		 "description": "Select the entity matcher file for overriding Fabric node URLs using pattern matching",
		 "type": "fileselector",
		 "fileExtensions": [
		  ".yaml"
		 ],
		 "visible": true,
		 "readonly": false,
		 "valid": true
		},
		"value": {
		 "filename": "local_entity_matchers.yaml",
		 "content": "data:application/octet-stream;base64,IwojIENvcHlyaWdodCBTZWN1cmVLZXkgVGVjaG5vbG9naWVzIEluYy4gQWxsIFJpZ2h0cyBSZXNlcnZlZC4KIwojIFNQRFgtTGljZW5zZS1JZGVudGlmaWVyOiBBcGFjaGUtMi4wCiMKIwojIFRoZSBuZXR3b3JrIGNvbm5lY3Rpb24gcHJvZmlsZSBwcm92aWRlcyBjbGllbnQgYXBwbGljYXRpb25zIHRoZSBpbmZvcm1hdGlvbiBhYm91dCB0aGUgdGFyZ2V0CiMgYmxvY2tjaGFpbiBuZXR3b3JrIHRoYXQgYXJlIG5lY2Vzc2FyeSBmb3IgdGhlIGFwcGxpY2F0aW9ucyB0byBpbnRlcmFjdCB3aXRoIGl0LiBUaGVzZSBhcmUgYWxsCiMga25vd2xlZGdlIHRoYXQgbXVzdCBiZSBhY3F1aXJlZCBmcm9tIG91dC1vZi1iYW5kIHNvdXJjZXMuIFRoaXMgZmlsZSBwcm92aWRlcyBzdWNoIGEgc291cmNlLgojCgojIEVudGl0eU1hdGNoZXJzIGVuYWJsZSBzdWJzdGl0dXRpb24gb2YgbmV0d29yayBob3N0bmFtZXMgd2l0aCBzdGF0aWMgY29uZmlndXJhdGlvbnMKICMgc28gdGhhdCBwcm9wZXJ0aWVzIGNhbiBiZSBtYXBwZWQuIFJlZ2V4IGNhbiBiZSB1c2VkIGZvciB0aGlzIHB1cnBvc2UKIyBVcmxTdWJzdGl0dXRpb25FeHAgY2FuIGJlIGVtcHR5IHdoaWNoIG1lYW5zIHRoZSBzYW1lIG5ldHdvcmsgaG9zdG5hbWUgd2lsbCBiZSB1c2VkCiMgVXJsU3Vic3RpdHV0aW9uRXhwIGNhbiBiZSBnaXZlbiBzYW1lIGFzIG1hcHBlZCBwZWVyIHVybCwgc28gdGhhdCBtYXBwZWQgcGVlciB1cmwgY2FuIGJlIHVzZWQKIyBVcmxTdWJzdGl0dXRpb25FeHAgY2FuIGhhdmUgZ29sYW5nIHJlZ2V4IG1hdGNoZXJzIGxpa2UgJHsxfS5leGFtcGxlLiR7Mn06JHszfSBmb3IgcGF0dGVybgogIyBsaWtlIHBlZXIwLm9yZzEuZXhhbXBsZS5jb206MTIzNCB3aGljaCBjb252ZXJ0cyBwZWVyMC5vcmcxLmV4YW1wbGUuY29tIHRvIHBlZXIwLm9yZzEuZXhhbXBsZS5jb206MTIzNAojIHNzbFRhcmdldE92ZXJyaWRlVXJsU3Vic3RpdHV0aW9uRXhwIGZvbGxvdyBpbiB0aGUgc2FtZSBsaW5lcyBhcwogIyBTdWJzdGl0dXRpb25FeHAgZm9yIHRoZSBmaWVsZHMgZ3ByY09wdGlvbnMuc3NsLXRhcmdldC1uYW1lLW92ZXJyaWRlIHJlc3BlY3RpdmVseQojIEluIGFueSBjYXNlIG1hcHBlZEhvc3QncyBjb25maWcgd2lsbCBiZSB1c2VkLCBzbyBtYXBwZWQgaG9zdCBjYW5ub3QgYmUgZW1wdHksIGlmIGVudGl0eU1hdGNoZXJzIGFyZSB1c2VkCmVudGl0eU1hdGNoZXJzOgogIHBlZXI6CiAgICAtIHBhdHRlcm46IHBlZXIwLm9yZzEuZXhhbXBsZS4oXHcrKQogICAgICB1cmxTdWJzdGl0dXRpb25FeHA6IGxvY2FsaG9zdDo3MDUxCiAgICAgIHNzbFRhcmdldE92ZXJyaWRlVXJsU3Vic3RpdHV0aW9uRXhwOiBwZWVyMC5vcmcxLmV4YW1wbGUuY29tCiAgICAgIG1hcHBlZEhvc3Q6IHBlZXIwLm9yZzEuZXhhbXBsZS5jb20KCiAgICAtIHBhdHRlcm46IHBlZXIxLm9yZzEuZXhhbXBsZS4oXHcrKQogICAgICB1cmxTdWJzdGl0dXRpb25FeHA6IGxvY2FsaG9zdDo4MDUxCiAgICAgIHNzbFRhcmdldE92ZXJyaWRlVXJsU3Vic3RpdHV0aW9uRXhwOiBwZWVyMS5vcmcxLmV4YW1wbGUuY29tCiAgICAgIG1hcHBlZEhvc3Q6IHBlZXIxLm9yZzEuZXhhbXBsZS5jb20KCiAgICAtIHBhdHRlcm46IHBlZXIwLm9yZzIuZXhhbXBsZS4oXHcrKQogICAgICB1cmxTdWJzdGl0dXRpb25FeHA6IGxvY2FsaG9zdDo5MDUxCiAgICAgIHNzbFRhcmdldE92ZXJyaWRlVXJsU3Vic3RpdHV0aW9uRXhwOiBwZWVyMC5vcmcyLmV4YW1wbGUuY29tCiAgICAgIG1hcHBlZEhvc3Q6IHBlZXIwLm9yZzIuZXhhbXBsZS5jb20KCiAgICAtIHBhdHRlcm46IHBlZXIxLm9yZzIuZXhhbXBsZS4oXHcrKQogICAgICB1cmxTdWJzdGl0dXRpb25FeHA6IGxvY2FsaG9zdDoxMDA1MQogICAgICBzc2xUYXJnZXRPdmVycmlkZVVybFN1YnN0aXR1dGlvbkV4cDogcGVlcjEub3JnMi5leGFtcGxlLmNvbQogICAgICBtYXBwZWRIb3N0OiBwZWVyMS5vcmcyLmV4YW1wbGUuY29tCgogICAgLSBwYXR0ZXJuOiAoXHcrKS5vcmcxLmV4YW1wbGUuKFx3Kyk6KFxkKykKICAgICAgdXJsU3Vic3RpdHV0aW9uRXhwOiBsb2NhbGhvc3Q6JHsyfQogICAgICBzc2xUYXJnZXRPdmVycmlkZVVybFN1YnN0aXR1dGlvbkV4cDogJHsxfS5vcmcxLmV4YW1wbGUuY29tCiAgICAgIG1hcHBlZEhvc3Q6ICR7MX0ub3JnMS5leGFtcGxlLmNvbQoKICAgIC0gcGF0dGVybjogKFx3Kykub3JnMi5leGFtcGxlLihcdyspOihcZCspCiAgICAgIHVybFN1YnN0aXR1dGlvbkV4cDogbG9jYWxob3N0OiR7Mn0KICAgICAgc3NsVGFyZ2V0T3ZlcnJpZGVVcmxTdWJzdGl0dXRpb25FeHA6ICR7MX0ub3JnMi5leGFtcGxlLmNvbQogICAgICBtYXBwZWRIb3N0OiAkezF9Lm9yZzIuZXhhbXBsZS5jb20KCiAgICAtIHBhdHRlcm46IChcdyspOjcwNTEKICAgICAgdXJsU3Vic3RpdHV0aW9uRXhwOiBsb2NhbGhvc3Q6NzA1MQogICAgICBzc2xUYXJnZXRPdmVycmlkZVVybFN1YnN0aXR1dGlvbkV4cDogcGVlcjAub3JnMS5leGFtcGxlLmNvbQogICAgICBtYXBwZWRIb3N0OiBwZWVyMC5vcmcxLmV4YW1wbGUuY29tCgogICAgLSBwYXR0ZXJuOiAoXHcrKTo4MDUxCiAgICAgIHVybFN1YnN0aXR1dGlvbkV4cDogbG9jYWxob3N0OjgwNTEKICAgICAgc3NsVGFyZ2V0T3ZlcnJpZGVVcmxTdWJzdGl0dXRpb25FeHA6IHBlZXIxLm9yZzEuZXhhbXBsZS5jb20KICAgICAgbWFwcGVkSG9zdDogcGVlcjEub3JnMS5leGFtcGxlLmNvbQoKICAgIC0gcGF0dGVybjogKFx3Kyk6OTA1MQogICAgICB1cmxTdWJzdGl0dXRpb25FeHA6IGxvY2FsaG9zdDo5MDUxCiAgICAgIHNzbFRhcmdldE92ZXJyaWRlVXJsU3Vic3RpdHV0aW9uRXhwOiBwZWVyMC5vcmcyLmV4YW1wbGUuY29tCiAgICAgIG1hcHBlZEhvc3Q6IHBlZXIwLm9yZzIuZXhhbXBsZS5jb20KCiAgICAtIHBhdHRlcm46IChcdyspOjEwMDUxCiAgICAgIHVybFN1YnN0aXR1dGlvbkV4cDogbG9jYWxob3N0OjEwMDUxCiAgICAgIHNzbFRhcmdldE92ZXJyaWRlVXJsU3Vic3RpdHV0aW9uRXhwOiBwZWVyMS5vcmcyLmV4YW1wbGUuY29tCiAgICAgIG1hcHBlZEhvc3Q6IHBlZXIxLm9yZzIuZXhhbXBsZS5jb20KCiAgb3JkZXJlcjoKCiAgICAtIHBhdHRlcm46IChcdyspLmV4YW1wbGUuKFx3KykKICAgICAgdXJsU3Vic3RpdHV0aW9uRXhwOiBsb2NhbGhvc3Q6NzA1MAogICAgICBzc2xUYXJnZXRPdmVycmlkZVVybFN1YnN0aXR1dGlvbkV4cDogb3JkZXJlci5leGFtcGxlLmNvbQogICAgICBtYXBwZWRIb3N0OiBvcmRlcmVyLmV4YW1wbGUuY29tCgogICAgLSBwYXR0ZXJuOiAoXHcrKS5leGFtcGxlLihcdyspOihcZCspCiAgICAgIHVybFN1YnN0aXR1dGlvbkV4cDogbG9jYWxob3N0OjcwNTAKICAgICAgc3NsVGFyZ2V0T3ZlcnJpZGVVcmxTdWJzdGl0dXRpb25FeHA6IG9yZGVyZXIuZXhhbXBsZS5jb20KICAgICAgbWFwcGVkSG9zdDogb3JkZXJlci5leGFtcGxlLmNvbQoKICBjZXJ0aWZpY2F0ZUF1dGhvcml0eToKICAgIC0gcGF0dGVybjogKFx3Kykub3JnMS5leGFtcGxlLihcdyspCiAgICAgIHVybFN1YnN0aXR1dGlvbkV4cDogaHR0cHM6Ly9sb2NhbGhvc3Q6NzA1NAogICAgICBzc2xUYXJnZXRPdmVycmlkZVVybFN1YnN0aXR1dGlvbkV4cDogY2Eub3JnMS5leGFtcGxlLmNvbQogICAgICBtYXBwZWRIb3N0OiBjYS5vcmcxLmV4YW1wbGUuY29tCgogICAgLSBwYXR0ZXJuOiAoXHcrKS5vcmcyLmV4YW1wbGUuKFx3KykKICAgICAgdXJsU3Vic3RpdHV0aW9uRXhwOiBodHRwczovL2xvY2FsaG9zdDo4MDU0CiAgICAgIHNzbFRhcmdldE92ZXJyaWRlVXJsU3Vic3RpdHV0aW9uRXhwOiBjYS5vcmcyLmV4YW1wbGUuY29tCiAgICAgIG1hcHBlZEhvc3Q6IGNhLm9yZzIuZXhhbXBsZS5jb20="
		}
	 },
	 {
		"name": "orgName",
		"type": "string",
		"required": true,
		"display": {
		 "name": "Client organization name",
		 "description": "Name of the organization that created the client user",
		 "visible": true,
		 "readonly": false,
		 "valid": true
		},
		"value": "org1"
	 },
	 {
		"name": "userName",
		"type": "string",
		"description": "Name of the client user",
		"required": true,
		"display": {
		 "name": "Client user name",
		 "visible": true,
		 "readonly": false,
		 "valid": true
		},
		"value": "User1"
	 }
	],
	"outputs": [],
	"inputs": [],
	"handler": {
	   "settings": []
	},
	"reply": [],
	"s3Prefix": "flogo",
	"key": "flogo/fabclient/connector/fabconnector/connector.json",
	"display": {
	   "category": "fabclient",
	   "description": "Fabric Connection",
	   "visible": true,
	   "smallIcon": "ic-fabconnector@2x.png",
	   "largeIcon": "ic-fabconnector@3x.png"
	},
	"actions": [
	   {
		"name": "Save Connector",
		"display": {
		 "readonly": false,
		 "valid": true,
		 "visible": true
		}
	   }
	],
	"feature": {},
	"propertyMap": {},
	"keyfield": "name",
	"isValid": true,
	"lastUpdatedTime": 1554219154480,
	"createdTime": 1554219154480,
	"user": "flogo",
	"subscriptionId": "flogo_sbsc",
	"connectorName": " ",
	"connectorDescription": " "
}`

func TestConnectorSettings(t *testing.T) {
	configs, err := GetSettings(connectionObject)
	require.NoError(t, err, "failed to extract settings")
	assert.NotEmpty(t, configs, "config should not be empty")
	for k, v := range configs {
		fmt.Printf("key %s value %T\n", k, v)
		if k == "entityMatcher" {
			content, err := ExtractFileContent(v)
			require.NoError(t, err, "failed to extract file content")
			fmt.Printf("entityMatcher: %s\n", string(content))
		}
	}
}
