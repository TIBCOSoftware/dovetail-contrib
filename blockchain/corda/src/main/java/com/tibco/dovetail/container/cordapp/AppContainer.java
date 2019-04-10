package com.tibco.dovetail.container.cordapp;

import com.tibco.dovetail.core.runtime.services.IContainerService;
import com.tibco.dovetail.core.runtime.services.IDataService;
import com.tibco.dovetail.core.runtime.services.IEventService;
import com.tibco.dovetail.core.runtime.services.ILogService;

import net.corda.core.crypto.Base58;
import net.corda.core.identity.AbstractParty;
import net.corda.core.node.ServiceHub;

public class AppContainer implements IContainerService {
	private static ServiceHub serviceHub = null;
	AppDataService dataService;
    AppEventService eventService = new AppEventService();
    AppLoggingService logService;
	AppFlow flowService;;
    
	public AppContainer(AppFlow flow) {
		serviceHub = flow.getServiceHub();
		this.flowService = flow;

		if(flow.getLogger() != null)
			logService = new AppLoggingService(flow.getLogger());
		else
			logService = new AppLoggingService(flow.getClass().getTypeName());
	
	
		dataService = new AppDataService(this.flowService.getServiceHub(), this.flowService.getTransactionBuilder());
	}
	
	@Deprecated
	public AppContainer(ServiceHub hub, AppFlow flow) {
		serviceHub = hub;
		this.flowService = flow;
		dataService = new AppDataService(hub, this.flowService.getTransactionBuilder());
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
	
	public ServiceHub getServiceHub() {
		return this.flowService.getServiceHub();
	}
	
	public AppFlow getFlowService() {
		return this.flowService;
	}
	
	 public static String partyToString(AbstractParty p) {
 		return Base58.encode(p.getOwningKey().getEncoded());
 }
 
	 public static AbstractParty partyFromString(String s) {
	 		return serviceHub.getIdentityService().partyFromKey(net.corda.core.crypto.Crypto.decodePublicKey(Base58.decode(s)));
	 }
}
