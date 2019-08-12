/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.example.iou;

import static org.junit.Assert.*;

import org.junit.Test;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.CordaContainer;
import com.tibco.dovetail.container.corda.CordaDataService;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.core.model.composer.HLCResource;
import com.tibco.dovetail.core.model.flow.FlowAppConfig;
import com.tibco.dovetail.core.runtime.compilers.App;
import com.tibco.dovetail.core.runtime.compilers.AppCompiler;
import com.tibco.dovetail.core.runtime.compilers.FlowCompiler;
import com.tibco.dovetail.core.runtime.engine.ContextImpl;
import com.tibco.dovetail.core.runtime.flow.ReplyData;
import com.tibco.dovetail.core.runtime.flow.TransactionFlow;
import com.tibco.dovetail.core.runtime.flow.TransactionFlows;
import com.tibco.dovetail.core.runtime.transaction.ITransactionService;
import com.tibco.dovetail.core.runtime.transaction.TxnInputAttribute;
import com.tibco.dovetail.core.runtime.trigger.ITrigger;
import com.tibco.dovetail.core.runtime.util.JsonUtil;

import junit.framework.Assert;
import net.corda.core.contracts.ContractState;
import net.corda.core.contracts.UniqueIdentifier;
import net.corda.core.identity.CordaX500Name;
import net.corda.core.identity.Party;
import net.corda.core.transactions.LedgerTransaction;
import net.corda.core.utilities.OpaqueBytes;
import net.corda.core.contracts.PartyAndReference;
import net.corda.finance.contracts.asset.Cash;
import net.corda.finance.contracts.asset.Cash.State;
import net.corda.testing.core.TestIdentity;

import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import static net.corda.finance.Currencies.DOLLARS;
import static net.corda.finance.Currencies.POUNDS;
//import static net.corda.testing.NodeTestUtils.ledger;


public class TestCordaFlowContract {
	IOU iou;
	IOUContractContract contract;

	CordaContainer ctnr;
	ContextImpl context;
	ITrigger trigger;
	
	Party bank = (new TestIdentity(new CordaX500Name("BigCorp", "New York", "GB"))).getParty();
		
	Party bob = (new TestIdentity(new CordaX500Name("bob", "New York", "GB"))).getParty();
		
	Party charlie = (new TestIdentity(new CordaX500Name("charlie", "New York", "GB"))).getParty();

	Party alice = (new TestIdentity(new CordaX500Name("alice", "New York", "GB"))).getParty();
	
	class MockIssueTxn implements ITransactionService {

		@Override
		public Map<String, Object> resolveTransactionInput(List<TxnInputAttribute> txnInputs) {
			Map<String, Object> context = new LinkedHashMap<String, Object>();
			DocumentContext doc = JsonUtil.getJsonParser().parse("{}");
			doc.put("$", "iou", CordaUtil.getInstance().toJsonObject((ContractState)iou).json());
		    doc.put("$", "transactionId", "issue");
		    doc.put("$", "timestamp", "abc");
		    context.put("transactionInput", doc);
		    context.put("containerServiceStub", ctnr);
		    return context;
		}

		@Override
		public boolean isTransactionSecuritySupported() {
			// TODO Auto-generated method stub
			return false;
		}

		@Override
		public String getTransactionName() {
			// TODO Auto-generated method stub
			return "com.example.iou.IssueIOU";
		}

		@Override
		public String getTransactionInitiator() {
			// TODO Auto-generated method stub
			return null;
		}

		@Override
		public String getInitiatorCertAttribute(String attr) {
			// TODO Auto-generated method stub
			return null;
		}
		
	}
	
	class MockTransferTxn implements ITransactionService {

		@Override
		public Map<String, Object> resolveTransactionInput(List<TxnInputAttribute> txnInputs) {
			Map<String, Object> context = new LinkedHashMap<String, Object>();
			DocumentContext doc = JsonUtil.getJsonParser().parse("{}");
		    doc.put("$", "iou", CordaUtil.getInstance().toJsonObject((ContractState)iou).json());
			doc.put("$", "newOwner",CordaUtil.getInstance().partyToString(charlie));
		    doc.put("$", "transactionId", "transfer");
		    doc.put("$", "timestamp", "abc");
		    System.out.println(doc.jsonString());
		    context.put("transactionInput", doc);
		    context.put("containerServiceStub", ctnr);
		    return context;
		}

		@Override
		public boolean isTransactionSecuritySupported() {
			// TODO Auto-generated method stub
			return false;
		}

		@Override
		public String getTransactionName() {
			// TODO Auto-generated method stub
			return "com.example.iou.TransferIOU";
		}

		@Override
		public String getTransactionInitiator() {
			// TODO Auto-generated method stub
			return null;
		}

		@Override
		public String getInitiatorCertAttribute(String attr) {
			// TODO Auto-generated method stub
			return null;
		}
		
	}
	
	class MockSettleTxn implements ITransactionService {
		String settleType;
		public MockSettleTxn(String type) {
			this.settleType = type;
		}
		@Override
		public Map<String, Object> resolveTransactionInput(List<TxnInputAttribute> txnInputs) {
			Map<String, Object> context = new LinkedHashMap<String, Object>();
			DocumentContext doc = JsonUtil.getJsonParser().parse("{}");
		    doc.put("$", "transactionId", this.settleType);
		    doc.put("$", "timestamp", "abc");
		    doc.put("$", "sendChangeTo", CordaUtil.getInstance().partyToString(alice));
		    doc.put("$", "sendPaymentTo", CordaUtil.getInstance().partyToString(bob));
		    
		    if(settleType.equals("single")) {
			    PartyAndReference issuer = new PartyAndReference(bank, OpaqueBytes.of("123".getBytes()));
	    			Cash.State payment1 = new Cash.State(issuer, DOLLARS(100), bob);
	    			DocumentContext paymentdoc = CordaUtil.getInstance().toJsonObject(Arrays.asList(payment1));
	    			doc.put("$", "iou", CordaUtil.getInstance().toJsonObject((ContractState)iou).json());
				doc.put("$", "funds", paymentdoc.json());
				doc.put("$", "payAmt", CordaUtil.getInstance().toJsonObject(DOLLARS(100)).json());
			    
		    } else if(settleType.equals("multiple")) {
			    	PartyAndReference issuer = new PartyAndReference(bank, OpaqueBytes.of("123".getBytes()));
	        		Cash.State payment1 = new Cash.State(issuer, DOLLARS(10), bob);
	        		Cash.State payment2 = new Cash.State(issuer, DOLLARS(50), bob);
	        		
	        		iou.setPaid(DOLLARS(10));
	        		doc.put("$", "iou", CordaUtil.getInstance().toJsonObject((ContractState)iou).json());
	    			doc.put("$", "funds", CordaUtil.getInstance().toJsonObject(Arrays.asList(payment1, payment2)).json());
	    			doc.put("$", "payAmt", CordaUtil.getInstance().toJsonObject(DOLLARS(60)).json());
	    			
		    } else if(settleType.equals("change")) {
			    	PartyAndReference issuer = new PartyAndReference(bank, OpaqueBytes.of("123".getBytes()));
	        		Cash.State payment1 = new Cash.State(issuer, DOLLARS(10), bob);
	        		Cash.State payment2 = new Cash.State(issuer, DOLLARS(50), bob);
	        		
	        		iou.setPaid(DOLLARS(45));
	        		doc.put("$", "iou", CordaUtil.getInstance().toJsonObject((ContractState)iou).json());
	    			doc.put("$", "funds", CordaUtil.getInstance().toJsonObject(Arrays.asList(payment1, payment2)).json());
	    			doc.put("$", "payAmt", CordaUtil.getInstance().toJsonObject(DOLLARS(55)).json());
		    } else if(settleType.equals("err")) {
			    	PartyAndReference issuer = new PartyAndReference(bank, OpaqueBytes.of("123".getBytes()));
	        		Cash.State payment1 = new Cash.State(issuer, DOLLARS(10), bob);
	        		Cash.State payment2 = new Cash.State(issuer, DOLLARS(50), bob);
	        		
	        		iou.setPaid(DOLLARS(45));
	        		doc.put("$", "iou", CordaUtil.getInstance().toJsonObject((ContractState)iou).json());
	    			doc.put("$", "funds", CordaUtil.getInstance().toJsonObject(Arrays.asList(payment1, payment2)).json());
	    			doc.put("$", "payAmt", CordaUtil.getInstance().toJsonObject(DOLLARS(60)).json());
	    			
		    }
		    context.put("transactionInput", doc);
		    context.put("containerServiceStub", ctnr);
		    return context;
		}
		
		@Override
		public boolean isTransactionSecuritySupported() {
			// TODO Auto-generated method stub
			return false;
		}

		@Override
		public String getTransactionName() {
			// TODO Auto-generated method stub
			return "com.example.iou.SettleIOU";
		}

		@Override
		public String getTransactionInitiator() {
			// TODO Auto-generated method stub
			return null;
		}

		@Override
		public String getInitiatorCertAttribute(String attr) {
			// TODO Auto-generated method stub
			return null;
		}
		
	}
	@org.junit.Before
	public void createIOU() {
		
		
		iou = new IOU(alice, bob, DOLLARS(100), DOLLARS(0), new UniqueIdentifier());
		contract = new IOUContractContract();
		
		List<ContractState> inputs = new ArrayList<ContractState>();
		inputs.add(iou);
		ctnr = new CordaContainer(inputs, "testlogger");
	
	    
	    InputStream txJson = this.getClass().getResourceAsStream("transactions.json");

		try {
			App app = AppCompiler.compileApp(txJson);
    	 		trigger = app.getTriggers().get("com.example.iou.IssueIOU");
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		
	}

	@Test
	public void testIssue() {
        	System.out.println("\ntestIssue...");
        try {
			ReplyData reply = trigger.invoke(ctnr, new MockIssueTxn());
             
            List<DocumentContext> out = ((CordaDataService)ctnr.getDataService()).getModifiedStates();
            
    			List<DocumentContext> expecteddocs = new ArrayList<DocumentContext>();
            expecteddocs.add(CordaUtil.getInstance().toJsonObject(iou));
            
            CordaUtil.getInstance().compare(expecteddocs, out);
            System.out.println(reply.getData());
            
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
     	System.out.println("\ntestIssue done");
	}
	
	@Test
	public void testTransfer() {
		System.out.println("\ntestTransfer...");
        
        try {
        		System.out.println("owner="+ CordaUtil.getInstance().partyToString(iou.getOwner()));
        		trigger.invoke(ctnr, new MockTransferTxn());
            
            List<DocumentContext> out = ((CordaDataService)ctnr.getDataService()).getModifiedStates();
            
    			List<DocumentContext> expecteddocs = new ArrayList<DocumentContext>();
    			iou.setOwner(charlie);
    			System.out.println("new owner="+ CordaUtil.getInstance().partyToString(iou.getOwner()));
            expecteddocs.add(CordaUtil.getInstance().toJsonObject(iou));
                
            CordaUtil.getInstance().compare(expecteddocs, out);
            
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
        System.out.println("\ntestTransfer...done");
	}
	
	@Test
	public void testSettleSinglePayment() {
		System.out.println("\ntestSettleSinglePayment....");
	
        try {
        		
        		trigger.invoke(ctnr, new MockSettleTxn("single"));
            List<DocumentContext> out = ((CordaDataService)ctnr.getDataService()).getModifiedStates();
            System.out.println("settle output state counts: " + out.size());
            
            List<DocumentContext> expecteddocs = new ArrayList<DocumentContext>();
            
            iou.setPaid(DOLLARS(100));
            DocumentContext expecteddoc = CordaUtil.getInstance().toJsonObject((ContractState)iou); 
            expecteddocs.add(expecteddoc);
            
            PartyAndReference issuer = new PartyAndReference(bank, OpaqueBytes.of("123".getBytes()));
            Cash.State payment1 = new Cash.State(issuer, DOLLARS(100), bob);
           
            expecteddocs.add(CordaUtil.getInstance().toJsonObject(payment1));
            
            CordaUtil.getInstance().compare(expecteddocs, out);
            
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
	
	@Test
	public void testSettleMultiplePayments() {
		System.out.println("\ntestSettleMultiplePayments....");
        try {
        		
			trigger.invoke(ctnr, new MockSettleTxn("multiple"));
           
            DocumentContext expecteddoc = CordaUtil.getInstance().toJsonObject((ContractState)iou);
            String expected = CordaUtil.getInstance().toJsonObject((ContractState)iou).jsonString();
            System.out.println("expected IOU=" + expected);
        
            List<DocumentContext> out = ((CordaDataService)ctnr.getDataService()).getModifiedStates();
            System.out.println("settle output state counts: " + out.size());
            out.forEach(doc -> {
	    			String actual = doc.jsonString();
	    			System.out.println(doc.jsonString());
            });
            
        		PartyAndReference issuer = new PartyAndReference(bank, OpaqueBytes.of("123".getBytes()));
            List<DocumentContext> expecteddocs = new ArrayList<DocumentContext>();
            Cash.State payment11 = new Cash.State(issuer, DOLLARS(60), bob);
            expecteddocs.add(CordaUtil.getInstance().toJsonObject(payment11));
            
            iou.setPaid(DOLLARS(70));
            expecteddocs.add(CordaUtil.getInstance().toJsonObject(iou));
            
            CordaUtil.getInstance().compare(expecteddocs, out);
          
            
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}

//	Test
	public void testSettleWithChange() {
		System.out.println("\ntestSettleWithChange....");
        try {
        		trigger.invoke(ctnr, new MockSettleTxn("change"));
            DocumentContext expecteddoc = CordaUtil.getInstance().toJsonObject((ContractState)iou);
            String expected = CordaUtil.getInstance().toJsonObject((ContractState)iou).jsonString();
        
            List<DocumentContext> out = ((CordaDataService)ctnr.getDataService()).getModifiedStates();
            System.out.println("settle output state counts: " + out.size());
            out.forEach(doc -> {
	    			String actual = doc.jsonString();
	    			System.out.println(doc.jsonString());
            });
            
            List<DocumentContext> expecteddocs = new ArrayList<DocumentContext>();
            PartyAndReference issuer = new PartyAndReference(bank, OpaqueBytes.of("123".getBytes()));
            Cash.State payment11 = new Cash.State(issuer, DOLLARS(55), alice);
            Cash.State payment22 = new Cash.State(issuer, DOLLARS(5), bob);
            expecteddocs.add(CordaUtil.getInstance().toJsonObject(payment11));
            expecteddocs.add(CordaUtil.getInstance().toJsonObject(payment22));
            
            iou.setPaid(DOLLARS(100));
            expecteddocs.add(CordaUtil.getInstance().toJsonObject(iou));
            
            CordaUtil.getInstance().compare(expecteddocs, out);
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
	
	@Test
	public void testSettleWithErr() {
		System.out.println("\ntestSettleWithErr...");

        try {
        		ReplyData reply = trigger.invoke(ctnr, new MockSettleTxn("err"));
            System.out.println("status=" + reply.getData());
            
		} catch (Exception e) {
			 
			assertTrue(e.getMessage().equals("incorrect payment"));
		}
	}
	
	@Test
	public void testCompare() {
		String a1 = "{\"owner\":\"GfHq2tTVk9z4eXgyFWjSLXiJwa9rNL8x2kfTqQJ38NfFx7DLD4hDyeRZbesJ\"," + 
				"  \"issuer\":\"GfHq2tTVk9z4eXgyFiQ7BQSoTC2EBtogVoZF72CgJxZ8BdzS3xpiGR8TRwtb\"," + 
				"  \"issuerRef\":\"AA==\", \"amt\":{\"quantity\":100, \"currency\":\"USD\"}}";
		
		String a2 = "{\"owner\":\"GfHq2tTVk9z4eXgyN5E5a5nEFC22yWtMykX9FoA1ph672xUjeJkK8HThbAoH\", " + 
				"  \"issuer\":\"GfHq2tTVk9z4eXgyFiQ7BQSoTC2EBtogVoZF72CgJxZ8BdzS3xpiGR8TRwtb\"," + 
				"  \"issuerRef\":\"AA==\", \"amt\":{\"quantity\":9900, \"currency\":\"USD\"}}";
		
		String a3 = "{\"issuer\":\"GfHq2tTVk9z4eXgyN5E5a5nEFC22yWtMykX9FoA1ph672xUjeJkK8HThbAoH\"," + 
				"  \"owner\":\"GfHq2tTVk9z4eXgyFWjSLXiJwa9rNL8x2kfTqQJ38NfFx7DLD4hDyeRZbesJ\"," + 
				"  \"amt\":{\"currency\":\"USD\", \"quantity\":100}, \"paid\":{\"currency\":\"USD\", \"quantity\":100}," + 
				"  \"linearId\":\"test#d2b1b064-c6a7-447d-94ec-23c8a9902580\"}";
		
		String f1 = "{\"owner\":\"GfHq2tTVk9z4eXgyFWjSLXiJwa9rNL8x2kfTqQJ38NfFx7DLD4hDyeRZbesJ\"," + 
				"  \"issuer\":\"GfHq2tTVk9z4eXgyFiQ7BQSoTC2EBtogVoZF72CgJxZ8BdzS3xpiGR8TRwtb\"," + 
				"  \"issuerRef\":\"AA==\", \"amt\":{\"quantity\":100, \"currency\":\"USD\"}}";
		
		String f2 = "{\"owner\":\"GfHq2tTVk9z4eXgyN5E5a5nEFC22yWtMykX9FoA1ph672xUjeJkK8HThbAoH\"," + 
				"  \"issuer\":\"GfHq2tTVk9z4eXgyFiQ7BQSoTC2EBtogVoZF72CgJxZ8BdzS3xpiGR8TRwtb\"," + 
				"  \"issuerRef\":\"AA==\", \"amt\":{\"quantity\":9900, \"currency\":\"USD\"}}";
		
		String f3 = "{\"issuer\":\"GfHq2tTVk9z4eXgyN5E5a5nEFC22yWtMykX9FoA1ph672xUjeJkK8HThbAoH\"," + 
				"  \"owner\":\"GfHq2tTVk9z4eXgyFWjSLXiJwa9rNL8x2kfTqQJ38NfFx7DLD4hDyeRZbesJ\"," + 
				"  \"amt\":{\"currency\":\"USD\", \"quantity\":100}, \"paid\":{\"currency\":\"USD\", \"quantity\":100}," + 
				"  \"linearId\":\"test#d2b1b064-c6a7-447d-94ec-23c8a9902580\"}";
		
		DocumentContext da1 = JsonUtil.getJsonParser().parse(a1);
		DocumentContext da2 = JsonUtil.getJsonParser().parse(a2);
		DocumentContext da3 = JsonUtil.getJsonParser().parse(a3);
		DocumentContext df1 = JsonUtil.getJsonParser().parse(f1);
		DocumentContext df2 = JsonUtil.getJsonParser().parse(f2);
		DocumentContext df3 = JsonUtil.getJsonParser().parse(f3);
		
		List<Map<String, Object>> av = new ArrayList<>();
        List<Map<String, Object>> rv = new ArrayList<>();

      //  ObjectMapper mapper = new ObjectMapper();
     
        		//av.add(mapper.readValue(actual.get(i).jsonString(), Map.class));
        		//rv.add(mapper.readValue(results.get(i).jsonString(), Map.class));
        		av.add(da1.json());
        		av.add(da2.json());
        		av.add(da3.json());
        		rv.add(df3.json());
        		rv.add(df2.json());
        		rv.add(df1.json());
        
        		System.out.println("a=f?" + av.containsAll(rv));
	}
}
