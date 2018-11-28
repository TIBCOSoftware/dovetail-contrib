/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package org.global.citizens.net;

import static net.corda.finance.Currencies.DOLLARS;

import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Map;

import org.junit.Test;

import org.global.citizens.net.*;
import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.CordaContainer;
import com.tibco.dovetail.container.corda.CordaDataService;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.core.runtime.compilers.FlowCompiler;
import com.tibco.dovetail.core.runtime.engine.ContextImpl;
import com.tibco.dovetail.core.runtime.engine.DovetailEngine;
import com.tibco.dovetail.core.runtime.flow.ReplyData;
import com.tibco.dovetail.core.runtime.flow.TransactionFlow;
import com.tibco.dovetail.core.runtime.flow.TransactionFlows;

import net.corda.core.contracts.ContractState;
import net.corda.core.contracts.UniqueIdentifier;
import smartcontract.trigger.transaction.ModelSchemaCompiler;
import smartcontract.trigger.transaction.model.composer.HLCResource;

import static net.corda.testing.TestConstants.getALICE;
import static net.corda.testing.TestConstants.getBOB;
import static net.corda.testing.TestConstants.getCHARLIE;

public class TestGC {
	TransactionFlows flows;
	/*
//	@org.junit.Before
	public void setupFlow() {
	
		ProjectPledgeContract contract = new ProjectPledgeContract();
		
		
	    InputStream txJson = contract.getTransactionJson();
        InputStream modelJson = contract.getSchemasJson();
        
        Map<String, HLCResource> schemas;
		try {
			schemas = ModelSchemaCompiler.compile(modelJson);
			flows = FlowCompiler.compile(txJson, schemas);
		} catch (IOException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		
	}
//	@Test
	public void testCreatePledge() {

        try {
        		System.out.println("\nCreateProjectPledge...");
        		ProjectPledge origin_pledge = new ProjectPledge("musicForChildren", "help kids to learn music", DOLLARS(10000), DOLLARS(0), DOLLARS(0), "GLOBALCITIZENREVIEW", getALICE(), getBOB(), null, new UniqueIdentifier("musicForChildren", null));
   
        		ContextImpl context = new ContextImpl();
        		context.addInput("pledgeId", "musicForChildren");
        		context.addInput("name", "musicForChildren");
        		context.addInput("description", "help kids to learn music");
        		context.addInput("fundsRequired", CordaUtil.toJsonObject(DOLLARS(10000)));
        		context.addInput("aidOrg", getALICE().toString());
        		context.addInput("globalCitizen", getBOB().toString());

        		context.addInput("containerServiceStub", context.getContainerService());
        		context.setContainerService(new CordaContainer(new ArrayList<ContractState>(), "gclogger"));
    
        		context.addInput("transactionId", "systxn1");
        		context.addInput("timestamp", "systimestamp");
			TransactionFlow flow = flows.getTransactionFlow("org.global.citizens.net.CreateProjectPledge");
			
			DovetailEngine engine = new DovetailEngine(flow);
            Object reply = engine.execute(context);
            
           
            List<DocumentContext> out = ((CordaDataService)context.getContainerService().getDataService()).getModifiedStates();
            
    			List<DocumentContext> expecteddocs = new ArrayList<DocumentContext>();
            expecteddocs.add(CordaUtil.toJsonObject(origin_pledge));
            
            CordaUtil.compare(expecteddocs, out);
            
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
	
	//@Test
	public void testRequestFunding() {

        try {
        		System.out.println("\nReuestFunding...");
        		ProjectPledge origin_pledge = new ProjectPledge("musicForChildren", "help kids to learn music", DOLLARS(10000), DOLLARS(0), DOLLARS(0), "GLOBALCITIZENREVIEW", getALICE(), getBOB(), null, new UniqueIdentifier("musicForChildren", null));
        		Funding funding = new Funding("", "REVIEW", DOLLARS(0), DOLLARS(0), DOLLARS(0), null, getCHARLIE());
        		ProjectPledge expected_pledge = new ProjectPledge("musicForChildren","help kids to learn music", DOLLARS(10000), DOLLARS(0), DOLLARS(0), "GOVERMENTREVIEW", getALICE(), getBOB(), Arrays.asList(funding), origin_pledge.getLinearId());
        		
        		ContextImpl context = new ContextImpl();
        		context.addInput("pledge", CordaUtil.toJsonObject(origin_pledge));
        		context.addInput("globalCitizen", getBOB().toString());
        		context.addInput("govOrg", Arrays.asList(getCHARLIE().toString()));
        		context.addInput("containerServiceStub", context.getContainerService());
        		context.setContainerService(new CordaContainer(new ArrayList<ContractState>(), "gclogger"));
    
        		context.addInput("transactionId", "systxn1");
        		context.addInput("timestamp", "systimestamp");
			TransactionFlow flow = flows.getTransactionFlow("org.global.citizens.net.RequestForFunding");
			
			DovetailEngine engine = new DovetailEngine(flow);
            Object reply = engine.execute(context);
            
           
            List<DocumentContext> out = ((CordaDataService)context.getContainerService().getDataService()).getModifiedStates();
            
    			List<DocumentContext> expecteddocs = new ArrayList<DocumentContext>();
            expecteddocs.add(CordaUtil.toJsonObject(expected_pledge));
            
            CordaUtil.compare(expecteddocs, out);
            
            
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
	
//	@Test
	public void testApproveFunding() {

        try {
        		System.out.println("\nApproveFunding...");
        		Funding funding = new Funding("", "REVIEW", DOLLARS(0), DOLLARS(0), DOLLARS(0), null, getCHARLIE());
        		ProjectPledge origin_pledge = new ProjectPledge("musicForChildren", "help kids to learn music", DOLLARS(10000), DOLLARS(0), DOLLARS(0), "GOVERMENTREVIEW", getALICE(), getBOB(), Arrays.asList(funding), new UniqueIdentifier("musicForChildren", null));
        		
        		Funding exfunding = new Funding("WEEKLY", "APPROVED", DOLLARS(5000), DOLLARS(0), DOLLARS(1000), null, getCHARLIE());
        		ProjectPledge expected_pledge = new ProjectPledge("musicForChildren",  "help kids to learn music", DOLLARS(10000), DOLLARS(5000), DOLLARS(0), "PROPOSALFUNDED", getALICE(), getBOB(), Arrays.asList(exfunding), origin_pledge.getLinearId());
        		
        		ContextImpl context = new ContextImpl();
        		context.addInput("pledge", CordaUtil.toJsonObject(origin_pledge));
        		context.addInput("govOrg", getCHARLIE().toString());
        		context.addInput("approvedFunding", CordaUtil.toJsonObject(DOLLARS(5000)));
        		context.addInput("fundsPerInstallment", CordaUtil.toJsonObject(DOLLARS(1000)));
        		context.addInput("fundingType", "WEEKLY");
        		context.addInput("containerServiceStub", context.getContainerService());
        		context.setContainerService(new CordaContainer(new ArrayList<ContractState>(), "gclogger"));
    
        		context.addInput("transactionId", "systxn1");
        		context.addInput("timestamp", "systimestamp");
			TransactionFlow flow = flows.getTransactionFlow("org.global.citizens.net.ApproveFunding");
			
			DovetailEngine engine = new DovetailEngine(flow);
            Object reply = engine.execute(context);
            
           
            List<DocumentContext> out = ((CordaDataService)context.getContainerService().getDataService()).getModifiedStates();
            
    			List<DocumentContext> expecteddocs = new ArrayList<DocumentContext>();
            expecteddocs.add(CordaUtil.toJsonObject(expected_pledge));
            
            CordaUtil.compare(expecteddocs, out);
            
            
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
	
//	@Test
	public void testTransferFunding() {

        try {
        		System.out.println("\nTransferFunding...");
        		
        		Funding funding = new Funding("WEEKLY", "APPROVED", DOLLARS(5000), DOLLARS(0), DOLLARS(1000), null, getCHARLIE());
        		ProjectPledge origin_pledge = new ProjectPledge("musicForChildren", "help kids to learn music", DOLLARS(10000), DOLLARS(5000), DOLLARS(0), "PROPOSALFUNDED", getALICE(), getBOB(), Arrays.asList(funding), new UniqueIdentifier("musicForChildren", null));
        		
        		Funding exfunding = new Funding("WEEKLY", "APPROVED", DOLLARS(5000), DOLLARS(1000), DOLLARS(1000), "2018-08-01", getCHARLIE());
        		ProjectPledge expected_pledge = new ProjectPledge("musicForChildren",  "help kids to learn music", DOLLARS(10000), DOLLARS(5000), DOLLARS(1000), "PROPOSALFUNDED", getALICE(), getBOB(), Arrays.asList(exfunding), origin_pledge.getLinearId());
        		
        		ContextImpl context = new ContextImpl();
        		context.addInput("pledge", CordaUtil.toJsonObject(origin_pledge));
        		context.addInput("govOrg", getCHARLIE().toString());
        		context.addInput("transferDate", "2018-08-01");

        		context.addInput("containerServiceStub", context.getContainerService());
        		context.setContainerService(new CordaContainer(new ArrayList<ContractState>(), "gclogger"));
    
        		context.addInput("transactionId", "systxn1");
        		context.addInput("timestamp", "systimestamp");
			TransactionFlow flow = flows.getTransactionFlow("org.global.citizens.net.TransferFunds");
			
			DovetailEngine engine = new DovetailEngine(flow);
            Object reply = engine.execute(context);
            
           
            List<DocumentContext> out = ((CordaDataService)context.getContainerService().getDataService()).getModifiedStates();
            
    			List<DocumentContext> expecteddocs = new ArrayList<DocumentContext>();
            expecteddocs.add(CordaUtil.toJsonObject(expected_pledge));
            
            CordaUtil.compare(expecteddocs, out);
            
            
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
	*/
}
