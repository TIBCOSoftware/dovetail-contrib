package com.tibco.dovetail.container.cordapp;
import java.util.LinkedHashMap;
import java.util.UUID;

import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.core.runtime.services.IContainerService;
import com.tibco.dovetail.core.runtime.services.IDataService;
import com.tibco.dovetail.core.runtime.services.IEventService;
import com.tibco.dovetail.core.runtime.services.ILogService;

import net.corda.core.node.ServiceHub;

public class AppContainer implements IContainerService {
//	private static ServiceHub serviceHub = null;
	AppDataService dataService;
    AppEventService eventService = new AppEventService();
    AppLoggingService logService;
    LinkedHashMap<String, Object> properties = new LinkedHashMap<String, Object>();
    
//	AppFlow flowService;;
    
	public AppContainer(AppFlow flow) {
		//serviceHub = flow.getServiceHub();
		CordaUtil.setServiceHub(flow.getServiceHub());
		
		this.properties.put("ServiceHub", flow.getServiceHub());
		this.properties.put("FlowService", flow);
		//this.flowService = flow;

		if(flow.getLogger() != null)
			logService = new AppLoggingService(flow.getLogger());
		else
			logService = new AppLoggingService(flow.getClass().getTypeName());
	
	
		dataService = new AppDataService(flow.getServiceHub(), flow.getTransactionBuilder(), flow.getRunId().getUuid());
	}
	
	@Deprecated
	public AppContainer(ServiceHub hub, AppFlow flow) {
		this.properties.put("ServiceHub", hub);
		this.properties.put("FlowService", flow);
		dataService = new AppDataService(hub, flow.getTransactionBuilder(), flow.getTransactionBuilder().getLockId());
		logService = new AppLoggingService("testing");
	}

	@Override
	public IDataService getDataService() {
		return this.dataService;
	}

	@Override
	public IEventService getEventService() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public ILogService getLogService() {
		return this.logService;
	}
	/*
	public ServiceHub getServiceHub() {
		return this.flowService.getServiceHub();
	}
	
	public AppFlow getFlowService() {
		return this.flowService;
	}
*/
	@Override
	public void addContainerProperty(String name, Object v) {
		this.properties.put(name, v);
	}

	@Override
	public Object getContainerProperty(String name) {
		return this.properties.get(name);
	}
}
