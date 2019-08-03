package com.tibco.dovetail.container.cordapp;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;

import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.core.runtime.services.IContainerService;
import com.tibco.dovetail.core.runtime.services.IDataService;
import com.tibco.dovetail.core.runtime.services.IEventService;
import com.tibco.dovetail.core.runtime.services.ILogService;

import net.corda.core.node.ServiceHub;

public class AppContainer implements IContainerService {

	public static final String TASK_SUBFLOW = "SUBFLOW";
	
	AppDataService dataService;
    AppEventService eventService = new AppEventService();
    AppLoggingService logService;
    LinkedHashMap<String, Object> properties = new LinkedHashMap<String, Object>();
    LinkedHashMap<String, List<Object>> tasks = new LinkedHashMap<String, List<Object>>();
    
    
	public AppContainer(AppFlow flow) {
		CordaUtil.setServiceHub(flow.getServiceHub());
		
		this.properties.put("ServiceHub", flow.getServiceHub());
		this.properties.put("FlowService", flow);

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
