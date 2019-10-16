/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package services

type ContainerService interface {
	GetDataService() DataService
	GetEventService() EventService
	GetLogService() LogService
}

type DataService interface {
	PutState(assetName string, assetKey string, assetValue map[string]interface{}, secondaryCompKeys interface{}) error
	GetState(assetName string, assetKey string, keyValue map[string]interface{}) ([]byte, error)
	DeleteState(assetName string, assetKey string, keyValue map[string]interface{}, secondaryCompKeys interface{}) ([]byte, error)
	LookupState(assetName string, assetKey string, lookupKey string, lkupKeyValue map[string]interface{}) ([][]byte, error)
	GetHistory(assetName string, assetKey string, keyValue map[string]interface{}) ([][]byte, error)
	QueryState(query string) ([][]byte, error)
}

type EventService interface {
	Publish(evtName, metadata string, evtPayload []byte) error
}

type LogService interface {
	Debug(msg string)
	Info(msg string)
	Warning(msg string)
	Error(errCode string, msg string, err error)
}
