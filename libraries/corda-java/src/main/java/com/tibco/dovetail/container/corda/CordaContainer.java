/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.dovetail.container.corda;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;

import com.tibco.dovetail.core.runtime.services.IContainerService;
import com.tibco.dovetail.core.runtime.services.IDataService;
import com.tibco.dovetail.core.runtime.services.IEventService;
import com.tibco.dovetail.core.runtime.services.ILogService;

import net.corda.core.contracts.ContractState;
import net.corda.core.flows.FlowLogicRefFactory;

public class CordaContainer implements IContainerService {
	public static final String TASK_SCHEDULEDACTIVITY = "SCHEDULED_ACTIVITY";
	CordaDataService dataService;
    CordaEventService eventService = new CordaEventService();
    CordaLoggingService logService;
    LinkedHashMap<String, Object> properties = new LinkedHashMap<String, Object>();
    LinkedHashMap<String, List<Object>> tasks = new LinkedHashMap<String, List<Object>>();
    
    FlowLogicRefFactory flowFactory;

	public CordaContainer(List<ContractState> inputs, String loggerName) {
		dataService = new CordaDataService(inputs);
		logService = new CordaLoggingService(loggerName);
	}
    
	public IDataService getDataService() {
		
		return dataService;
	}

	public IEventService getEventService() {
		return eventService;
	}

	public ILogService getLogService() {
		 return logService;
	}

	public void addContainerProperty(String name, Object v) {
		this.properties.put(name, v);
	}

	@Override
	public Object getContainerProperty(String name) {
		return this.properties.get(name);
	}

	@Override
	public void addContainerAsyncTask(String name, Object v) {
		List<Object> values = tasks.get(name);
		if(values == null) {
			values = new ArrayList<Object>();
			tasks.put(name, values);
		}
		values.add(v);
	}
	
	public LinkedHashMap<String, List<Object>> getContainerAsyncTasks() {
		return this.tasks;
	}
	
	public List<Object> getContainerAsyncTasks(String tasktype) {
		return this.tasks.get(tasktype);
	}
}
