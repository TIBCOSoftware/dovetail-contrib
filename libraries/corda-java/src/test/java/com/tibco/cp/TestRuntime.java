/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.cp;
import static org.junit.Assert.*;

import java.io.IOException;
import java.io.InputStream;
import java.util.List;
import java.util.Set;

import org.junit.Test;

import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jayway.jsonpath.DocumentContext;
import com.jayway.jsonpath.JsonPath;
import com.tibco.dovetail.core.model.flow.FlowAppConfig;
import com.tibco.dovetail.core.model.flow.HandlerConfig;
import com.tibco.dovetail.core.model.flow.Resources;
import com.tibco.dovetail.core.runtime.compilers.FlowCompiler;
import com.tibco.dovetail.core.runtime.engine.ContextImpl;
import com.tibco.dovetail.core.runtime.engine.FlowEngine;
import com.tibco.dovetail.core.runtime.flow.ReplyData;
import com.tibco.dovetail.core.runtime.flow.TransactionFlow;
import com.tibco.dovetail.core.runtime.engine.Scope;
import com.tibco.dovetail.core.runtime.services.IContainerService;
import com.tibco.dovetail.core.runtime.services.IDataService;
import com.tibco.dovetail.core.runtime.services.IEventService;
import com.tibco.dovetail.core.runtime.services.ILogService;

import junit.framework.Assert;

public class TestRuntime {

	//@Test
	public void testShcemaCompiler() throws Exception {
		ObjectMapper mapper = new ObjectMapper();
		InputStream in = this.getClass().getResourceAsStream("transactions.json");
		
			FlowAppConfig app = FlowAppConfig.parseModel(in);
			System.out.println(mapper.writeValueAsString(app));
			assertTrue(app.getTriggers().length == 1);
			
			HandlerConfig[] handlerConfigs = app.getTriggers()[0].getHandlers();
			for(int j=0; j<handlerConfigs.length; j++) {
				String txnName = String.valueOf(handlerConfigs[j].getSetting("transaction"));
				//Resources r = handlerConfigs[j].getFlow();
	
	             TransactionFlow flow = FlowCompiler.compile(handlerConfigs[j]);
			}
		
		
	}
	
	//@Test 
	public void testIterator() throws Exception {
		InputStream in = this.getClass().getResourceAsStream("iterator.json");
		
		FlowAppConfig app = FlowAppConfig.parseModel(in);
		
		FlowEngine e = new FlowEngine(FlowCompiler.compile(app.getTriggers()[0].getHandlers()[0]));
		
		ContextImpl context = new ContextImpl();
		context.setContainerService(new MockContainer());
		
		String auditjson = "[{\"user_txn_id\": \"id1\", \"hash_type\":\"rsa\", \"hash_value\":\"abcd\", \"data\": \"data1\"}, {\"user_txn_id\": \"id2\", \"hash_type\":\"rsa\", \"hash_value\":\"abcd\", \"data\": \"data2\"}]";
		DocumentContext doc = JsonPath.parse(auditjson);
		context.addInput("records", doc);
		context.addInput("transactionId", "first");
		context.addInput("timestamp", "timestamp");
		Scope scope = new Scope();
		ReplyData reply = e.execute(context, scope);
		Assert.assertEquals("Success", reply.getData());
	}
	
	public class MockContainer implements IContainerService {

		@Override
		public IDataService getDataService() {
			return new MockDataService();
		}

		@Override
		public IEventService getEventService() {
			return new MockEventService();
		}

		@Override
		public ILogService getLogService() {
			return new MockLogService();
		}

		@Override
		public void addContainerAsyncTask(String name, Object v){
			// Do nothing
		}

		@Override
		public Object getContainerProperty(String name){
			return "Property";
		}
		
	}
	
	public class MockDataService<T, R> implements IDataService<T, R> {


		@Override
		public R getState(String assetName, String assetKey, T keyValue) {
			// TODO Auto-generated method stub
			return null;
		}

		@Override
		public List<R> getStates(String assetName, String assetKey, Set<T> keyValue){
			return null;
		}

		@Override
		public R putState(String assetName, String assetKey, R assetValue){
			return assetValue;
		}

		@Override
		public R deleteState(String assetName, String assetKey, T keyValue) {
			// TODO Auto-generated method stub
			return null;
		}

		@Override
		public boolean processPayment(DocumentContext assetValue){
			return true;
		}

		@Override
		public List<R> lookupState(String assetName, String assetKey, T keyValue) {
			// TODO Auto-generated method stub
			return null;
		}

		@Override
		public List<R> getHistory(String assetName, String assetKey, T keyValue) {
			// TODO Auto-generated method stub
			return null;
		}

		@Override
		public List<DocumentContext> queryState(Object query) {
			// TODO Auto-generated method stub
			return null;
		}
	}
	
	public class MockEventService implements IEventService {

		@Override
		public void publish(String evtName, String metadata, String evtPayload) {
			// TODO Auto-generated method stub
			
		}
	}
	
	public class MockLogService implements ILogService {

		@Override
		public void debug(String msg) {
			// TODO Auto-generated method stub
			
		}

		@Override
		public void info(String msg) {
			// TODO Auto-generated method stub
			
		}

		@Override
		public void warning(String msg) {
			// TODO Auto-generated method stub
			
		}

		@Override
		public void error(String errCode, String msg, Throwable err) {
			// TODO Auto-generated method stub
			
		}
		
	}
	
}
