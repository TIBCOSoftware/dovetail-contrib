/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.dovetail.container.corda;

import java.util.List;

import com.tibco.dovetail.core.runtime.services.IContainerService;
import com.tibco.dovetail.core.runtime.services.IDataService;
import com.tibco.dovetail.core.runtime.services.IEventService;
import com.tibco.dovetail.core.runtime.services.ILogService;

import net.corda.core.contracts.ContractState;

public class CordaContainer implements IContainerService {
	CordaDataService dataService;
    CordaEventService eventService = new CordaEventService();
    CordaLoggingService logService;

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
    
}
