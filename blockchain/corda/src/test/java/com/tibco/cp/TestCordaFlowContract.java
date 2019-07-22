/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.cp;

import static org.junit.Assert.*;

import org.junit.Test;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.CordaContainer;
import com.tibco.dovetail.container.corda.CordaDataService;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.core.model.composer.HLCResource;
import com.tibco.dovetail.core.model.flow.FlowAppConfig;
import com.tibco.dovetail.core.runtime.compilers.FlowCompiler;
import com.tibco.dovetail.core.runtime.engine.ContextImpl;
import com.tibco.dovetail.core.runtime.engine.DovetailEngine;
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
import smartcontract.trigger.transaction.ModelSchemaCompiler;

import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import javax.print.Doc;

import static net.corda.finance.Currencies.DOLLARS;
import static net.corda.finance.Currencies.POUNDS;
import static net.corda.finance.Currencies.issuedBy;
//import static net.corda.testing.NodeTestUtils.ledger;


public class TestCordaFlowContract {
	com.tibco.cp.IOU3 iou;
	com.tibco.cp.IOU3Contract contract;

	CordaContainer ctnr;
	ContextImpl context;
	ITrigger trigger;
	
	static Party getMEGA_CORP() {
		TestIdentity mega = new TestIdentity(new CordaX500Name("BigCorp", "New York", "GB"));
		return mega.getParty();
	}
	
	static Party getBOB() {
		TestIdentity mega = new TestIdentity(new CordaX500Name("bob", "New York", "GB"));
		return mega.getParty();
	}
	
	static Party getCHARLIE() {
		TestIdentity mega = new TestIdentity(new CordaX500Name("charlie", "New York", "GB"));
		return mega.getParty();
	}
	
	static Party getALICE() {
		TestIdentity mega = new TestIdentity(new CordaX500Name("alice", "New York", "GB"));
		return mega.getParty();
	}
	
	class MockIssueTxn implements ITransactionService {

		@Override
		public Map<String, Object> resolveTransactionInput(List<TxnInputAttribute> txnInputs) {
			Map<String, Object> context = new LinkedHashMap<String, Object>();
			DocumentContext doc = JsonUtil.getJsonParser().parse("{}");
			doc.put("$", "iou", CordaUtil.toJsonObject((ContractState)iou).json());
		    doc.put("$", "transactionId", "issue");
		    doc.put("$", "timestamp", "abc");
		    context.put("transactionInput", doc);
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
			return "com.tibco.cp.IssueIOU";
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
		    context.put("iou", CordaUtil.toJsonObject((ContractState)iou));
			context.put("newLender",CordaUtil.serialize(getCHARLIE()));
		    context.put("transactionId", "transfer");
		    context.put("timestamp", "abc");
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
			return "com.tibco.cp.TransferIOU";
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
		    context.put("transactionId", this.settleType);
		    context.put("timestamp", "abc");
		    
		    if(settleType.equals("single")) {
			    PartyAndReference issuer = new PartyAndReference(getMEGA_CORP(), OpaqueBytes.of("123".getBytes()));
	    			Cash.State payment1 = new Cash.State(issuer, DOLLARS(10), getBOB());
	    			DocumentContext paymentdoc = CordaUtil.toJsonObject(Arrays.asList(payment1));
	    			 context.put("iou", CordaUtil.toJsonObject((ContractState)iou));
				context.put("payments", paymentdoc);
			    return context;
		    } else if(settleType.equals("multiple")) {
			    	PartyAndReference issuer = new PartyAndReference(getMEGA_CORP(), OpaqueBytes.of("123".getBytes()));
	        		Cash.State payment1 = new Cash.State(issuer, DOLLARS(10), getBOB());
	        		Cash.State payment2 = new Cash.State(issuer, DOLLARS(50), getBOB());
	        		
	        		iou.setPaid(DOLLARS(10));
	        		context.put("iou", CordaUtil.toJsonObject((ContractState)iou));
	    			context.put("payments", CordaUtil.toJsonObject(Arrays.asList(payment1, payment2)));
	    			System.out.println("expected Cash=" + CordaUtil.serialize(payment1));
		    } else if(settleType.equals("change")) {
			    	PartyAndReference issuer = new PartyAndReference(getMEGA_CORP(), OpaqueBytes.of("123".getBytes()));
	        		Cash.State payment1 = new Cash.State(issuer, DOLLARS(10), getBOB());
	        		Cash.State payment2 = new Cash.State(issuer, DOLLARS(50), getBOB());
	        		
	        		iou.setPaid(DOLLARS(45));
	        		context.put("iou", CordaUtil.toJsonObject((ContractState)iou));
	    			context.put("payments", CordaUtil.toJsonObject(Arrays.asList(payment1, payment2)));
		    } else if(settleType.equals("mixed")) {
			    	PartyAndReference issuer = new PartyAndReference(getMEGA_CORP(), OpaqueBytes.of("123".getBytes()));
	        		Cash.State payment1 = new Cash.State(issuer, DOLLARS(10), getBOB());
	        		Cash.State payment2 = new Cash.State(issuer, POUNDS(50), getBOB());
	        		
	        		iou.setPaid(DOLLARS(45));
	        		context.put("iou", CordaUtil.toJsonObject((ContractState)iou));
	    			context.put("payments", CordaUtil.toJsonObject(Arrays.asList(payment1, payment2)));
	    			System.out.println("expected Cash=" + CordaUtil.serialize(payment1));
		    }
		    
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
			return "com.tibco.cp.SettleIOU";
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
		
		
		iou = new IOU3(getALICE(), getBOB(), DOLLARS(100), DOLLARS(0), new UniqueIdentifier());
		contract = new IOU3Contract();
		
		List<ContractState> inputs = new ArrayList<ContractState>();
		inputs.add(iou);
		ctnr = new CordaContainer(inputs, "testlogger");
	
	    
	    InputStream txJson = contract.getTransactionJson();

		try {
			FlowAppConfig app = FlowAppConfig.parseModel(txJson);
    	 		DovetailEngine engine = new DovetailEngine(app);
    	 		trigger = engine.getTrigger("IssueIOU");
		} catch (IOException e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		
	}

	@Test
	public void testIssue() {
        	System.out.println("\ntestIssue...");
        try {
			trigger.invoke(ctnr, new MockIssueTxn());
             
            List<DocumentContext> out = ((CordaDataService)ctnr.getDataService()).getModifiedStates();
            
    			List<DocumentContext> expecteddocs = new ArrayList<DocumentContext>();
            expecteddocs.add(CordaUtil.toJsonObject(iou));
            
            CordaUtil.compare(expecteddocs, out);
            
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
        		
        		trigger.invoke(ctnr, new MockTransferTxn());
            
            List<DocumentContext> out = ((CordaDataService)ctnr.getDataService()).getModifiedStates();
            
    			List<DocumentContext> expecteddocs = new ArrayList<DocumentContext>();
    			iou.setLender(getCHARLIE());
            expecteddocs.add(CordaUtil.toJsonObject(iou));
                
            CordaUtil.compare(expecteddocs, out);
            
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
        System.out.println("\ntestTransfer...done");
	}
	
//	@Test
	public void testSettleSinglePayment() {
		System.out.println("\ntestSettleSinglePayment....");
		context.addInput("transactionId", "settle");
        try {
        		
        		trigger.invoke(ctnr, new MockSettleTxn("single"));
            List<DocumentContext> out = ((CordaDataService)ctnr.getDataService()).getModifiedStates();
            System.out.println("settle output state counts: " + out.size());
            
            List<DocumentContext> expecteddocs = new ArrayList<DocumentContext>();
            
            iou.setPaid(DOLLARS(10));
            DocumentContext expecteddoc = CordaUtil.toJsonObject((ContractState)iou); 
            expecteddocs.add(expecteddoc);
            
            PartyAndReference issuer = new PartyAndReference(getMEGA_CORP(), OpaqueBytes.of("123".getBytes()));
            Cash.State payment1 = new Cash.State(issuer, DOLLARS(10), getBOB());
            payment1 = (State) payment1.withNewOwner(getALICE()).getOwnableState();
            expecteddocs.add(CordaUtil.toJsonObject(payment1));
            
            CordaUtil.compare(expecteddocs, out);
            
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
	
//	@Test
	public void testSettleMultiplePayments() {
		System.out.println("\ntestSettleMultiplePayments....");
        try {
        		
			trigger.invoke(ctnr, new MockSettleTxn("multiple"));
           
            DocumentContext expecteddoc = CordaUtil.toJsonObject((ContractState)iou);
            String expected = CordaUtil.toJsonObject((ContractState)iou).jsonString();
            System.out.println("expected IOU=" + expected);
        
            List<DocumentContext> out = ((CordaDataService)ctnr.getDataService()).getModifiedStates();
            System.out.println("settle output state counts: " + out.size());
            out.forEach(doc -> {
	    			String actual = doc.jsonString();
	    			System.out.println(doc.jsonString());
            });
            
        		PartyAndReference issuer = new PartyAndReference(getMEGA_CORP(), OpaqueBytes.of("123".getBytes()));
            List<DocumentContext> expecteddocs = new ArrayList<DocumentContext>();
            Cash.State payment11 = new Cash.State(issuer, DOLLARS(60), getALICE());
            expecteddocs.add(CordaUtil.toJsonObject(payment11));
            
            iou.setPaid(DOLLARS(70));
            expecteddocs.add(CordaUtil.toJsonObject(iou));
            
            CordaUtil.compare(expecteddocs, out);
          
            
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}

//	@Test
	public void testSettleWithChange() {
		System.out.println("\ntestSettleWithChange....");
        try {
        		trigger.invoke(ctnr, new MockSettleTxn("change"));
            DocumentContext expecteddoc = CordaUtil.toJsonObject((ContractState)iou);
            String expected = CordaUtil.toJsonObject((ContractState)iou).jsonString();
        
            List<DocumentContext> out = ((CordaDataService)ctnr.getDataService()).getModifiedStates();
            System.out.println("settle output state counts: " + out.size());
            out.forEach(doc -> {
	    			String actual = doc.jsonString();
	    			System.out.println(doc.jsonString());
            });
            
            List<DocumentContext> expecteddocs = new ArrayList<DocumentContext>();
            PartyAndReference issuer = new PartyAndReference(getMEGA_CORP(), OpaqueBytes.of("123".getBytes()));
            Cash.State payment11 = new Cash.State(issuer, DOLLARS(55), getALICE());
            Cash.State payment22 = new Cash.State(issuer, DOLLARS(5), getBOB());
            expecteddocs.add(CordaUtil.toJsonObject(payment11));
            expecteddocs.add(CordaUtil.toJsonObject(payment22));
            
            iou.setPaid(DOLLARS(100));
            expecteddocs.add(CordaUtil.toJsonObject(iou));
            
            CordaUtil.compare(expecteddocs, out);
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
	}
	
//	@Test
	public void testSettleMixedCurrency() {
		System.out.println("\ntestSettleMixedCurrency...");
        try {
        		trigger.invoke(ctnr, new MockSettleTxn("mixed"));
            DocumentContext expecteddoc = CordaUtil.toJsonObject((ContractState)iou);
            String expected = CordaUtil.toJsonObject((ContractState)iou).jsonString();
            System.out.println("expected IOU=" + expected);
        
            List<DocumentContext> out = ((CordaDataService)ctnr.getDataService()).getModifiedStates();
            System.out.println("settle output state counts: " + out.size());
            out.forEach(doc -> {
	    			String actual = doc.jsonString();
	    			System.out.println(doc.jsonString());
            });
            
		} catch (Exception e) {
			assertTrue(e.getMessage().equals("payments must have the same currecy as IOU"));
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
