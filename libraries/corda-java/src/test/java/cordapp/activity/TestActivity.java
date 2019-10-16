package cordapp.activity;

import static net.corda.finance.Currencies.DOLLARS;
import static org.junit.Assert.*;

import java.io.IOException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import org.junit.Test;

import com.example.cp.IOU;
//import com.example.iou.IOU;
import com.fasterxml.jackson.core.JsonProcessingException;

import com.google.common.collect.Lists;
import com.jayway.jsonpath.DocumentContext;
import com.jayway.jsonpath.ParseContext;
import com.tibco.dovetail.container.corda.CordaContainer;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.container.cordapp.AppContainer;
import com.tibco.dovetail.container.cordapp.AppFlow;
import com.tibco.dovetail.container.cordapp.AppUtil;
import com.tibco.dovetail.core.model.metadata.MetadataParser;
import com.tibco.dovetail.core.runtime.engine.ContextImpl;
import com.tibco.dovetail.core.runtime.util.JsonUtil;

import dovetail_cordapp.activity.txnbuilder.txnbuilder;
import dovetail_cordapp.activity.wallet.wallet;
import dovetail_general.activity.sum.sum;
import junit.framework.Assert;
import net.corda.core.contracts.ContractState;
import net.corda.core.contracts.StateAndRef;
import net.corda.core.contracts.TransactionState;
import net.corda.core.contracts.UniqueIdentifier;
import net.corda.core.flows.FlowException;
import net.corda.core.identity.CordaX500Name;
import net.corda.core.transactions.SignedTransaction;
import net.corda.testing.core.TestIdentity;
import net.corda.testing.node.MockServices;
import dovetail_ledger.activity.payment.payment;

import net.corda.testing.node.MockNetwork;
import net.corda.testing.node.MockNetworkParameters;
import net.corda.testing.node.StartedMockNode;
import org.junit.After;
import org.junit.Before;

import static java.util.Collections.singletonList;
import static net.corda.testing.node.TestCordapp.findCordapp;

public class TestActivity {
	
	// private final MockNetwork mockNet = new MockNetwork(new MockNetworkParameters(singletonList(findCordapp("com.example.cp"))));
	// private StartedMockNode nodeA;
	// private StartedMockNode nodeB;
	    
	 @Before
	    public void setUp() {
	  //      nodeA = mockNet.createNode();
	        // We can optionally give the node a name.
	      //  nodeB = mockNet.createNode(new CordaX500Name("Bank B", "London", "GB"));
	        
	    }
	 
	 @After
    public void cleanUp() {
      //  mockNet.stopNodes();
    }
	 
	@Test
	public void testTxbuilder() throws IOException {
		 TestIdentity bob, alice;
			
			
		bob = new TestIdentity(new CordaX500Name("bob", "New York", "GB"));
		alice = new TestIdentity(new CordaX500Name("alice", "New York", "GB"));
		MockServices mock = new MockServices(
		        bob,
		        alice
		);
		AppUtil.setServiceHub(mock);
		
		txnbuilder builder = new txnbuilder();
		
		ContextImpl context = new ContextImpl();
		
		context.addInput("command", "com.example.cp.IssueIOU");
		String schema = "{\"schema\":\"http://json-schema.org/draft-04/schema#\",\"type\":\"object\",\"properties\":{\"iou\":{\"type\":\"object\",\"properties\":{\"linearId\":{\"type\":\"string\"},\"issuer\":{\"type\":\"string\"},\"owner\":{\"type\":\"string\"},\"amt\":{\"type\":\"object\",\"properties\":{\"currency\":{\"type\":\"string\"},\"quantity\":{\"type\":\"number\"}}},\"paid\":{\"type\":\"object\",\"properties\":{\"currency\":{\"type\":\"string\"},\"quantity\":{\"type\":\"number\"}}},\"createDt\":{\"type\":\"string\",\"format\":\"date-time\"}}}},\"description\":\"{\\\"metadata\\\":{\\\"type\\\":\\\"Transaction\\\",\\\"parent\\\":\\\"\\\",\\\"actors\\\":[\\\"issuer|\\\"],\\\"asset\\\":\\\"com.example.cp.IOU\\\",\\\"timewindow\\\":{}},\\\"attributes\\\":[{\\\"name\\\":\\\"iou\\\",\\\"type\\\":\\\"com.example.cp.IOU\\\",\\\"isRef\\\":false,\\\"isArray\\\":false,\\\"isReferenceData\\\":false,\\\"isAsset\\\":true,\\\"isParticipant\\\":false}]}\"}";
                
		context.addSetting("input_metadata", MetadataParser.parseSingleSchema(schema));
		
		//AppContainer ctnr = new AppContainer(mock);
		MockFlow flowservice = new MockFlow(true);
		
		AppContainer ctnr = new AppContainer(mock, flowservice);
		context.setContainerService(ctnr);
		
		String issuer = CordaUtil.getInstance().partyToString(bob.getParty());
		DocumentContext doc = JsonUtil.getJsonParser().parse("{}");
	//	DocumentContext iou = JsonUtil.getJsonParser().parse("{\"issuer\":\"" + issuer + "\", \"owner\":\"" + issuer + "\", \"amt\":{\"currency\":\"USD\", \"quantity\":100}, \"linearId\":\"myextid\"}");
		IOU iou = new IOU(java.time.Instant.now(), bob.getParty(), bob.getParty(), DOLLARS(100), DOLLARS(0), new UniqueIdentifier());
		Map iouvalue = CordaUtil.getInstance().toJsonObject((ContractState)iou).json();
	
		doc.put("$", "iou", iouvalue);

		context.addInput("input", doc);
		builder.eval(context);
		AppFlow txservice = (AppFlow) context.getContainerService().getContainerProperty("FlowService");
		assertEquals(1, txservice.getOutputStates().size());
		System.out.println("input=" + CordaUtil.getInstance().toJsonObject(iou).jsonString());
		System.out.println("output=" + CordaUtil.getInstance().toJsonObject(txservice.getOutputStates().get(0)).jsonString());
		assertEquals(iou.getParticipants(), txservice.getOutputStates().get(0).getParticipants());
		

		
	}
	
	//@Test
	public void testwallet() {
		MockServices mock;
		TestIdentity bob;
		ContextImpl context = new ContextImpl();
		
		bob = new TestIdentity(new CordaX500Name("BigCorp", "New York", "GB"));
		mock = new MockServices(bob);
		AppUtil.setServiceHub(mock);
		AppContainer ctnr = new AppContainer(mock, new MockFlow(true));
		context.setContainerService(ctnr);
		
		DocumentContext doc = JsonUtil.getJsonParser().parse("{\"amt\":{\"currency\":\"USD\", quantity:100}, \"issuers\":[\"" + CordaUtil.getInstance().partyToString(bob.getParty()) + "\"]}");
		
		
		
		context.addInput("operation", "Retrieve Funds");
		context.addInput("input", doc);
		
		wallet wal = new wallet();
		wal.eval(context);
		
		
	}
	
	@Test
	public void testsum() {
			
			ContextImpl context = new ContextImpl();
			
			ParseContext parser = JsonUtil.getJsonParser();
			DocumentContext doc = JsonUtil.getJsonParser().parse("{\"dataset\":[{\"field\":100.00},{\"field\":200}]}");
			
			context.addInput("scale", "2");
			context.addInput("rounding", "HALF_EVEN");
			context.addInput("input", doc);
			
			sum sumtest = new dovetail_general.activity.sum.sum();
			sumtest.eval(context);
			LinkedHashMap out = ((DocumentContext) context.getOutput("output")).json();
			assertEquals("300.00", out.get("result").toString());
		}
	
	@Test
	public void testdistinct() {
			
			ContextImpl context = new ContextImpl();
		
			DocumentContext doc = JsonUtil.getJsonParser().parse("{\"dataset\":[{\"field\":\"a\"},{\"field\":\"b\"}]}");
			
			
			context.addInput("input", doc);
			context.addInput("operation", "DISTINCT");
			
			dovetail_general.activity.collection.collection col = new dovetail_general.activity.collection.collection();
			col.eval(context);
			LinkedHashMap out = ((DocumentContext) context.getOutput("output")).json();
			assertEquals(2, out.get("size"));
		}
	
		@Test
		public void testpayment() {
			MockServices mock;
			TestIdentity bob;
			ContextImpl context = new ContextImpl();
			
			bob = new TestIdentity(new CordaX500Name("BigCorp", "New York", "GB"));
			TestIdentity charlie = new TestIdentity(new CordaX500Name("Charlie", "New York", "GB"));
			TestIdentity alice = new TestIdentity(new CordaX500Name("alice", "New York", "GB"));
			TestIdentity john = new TestIdentity(new CordaX500Name("alice", "New York", "GB"));
			mock = new MockServices(bob);
			
			String json = "{\"funds\":[{\"amt\":{\"currency\":\"USD\", \"quantity\":100}, \"issuer\":\"" + CordaUtil.getInstance().partyToString(bob.getParty()) + "\",\"issuerRef\":\"100\",\"owner\":\"" + CordaUtil.getInstance().partyToString(charlie.getParty()) + "\"},{\"amt\":{\"currency\":\"USD\", \"quantity\":100}, \"issuer\":\"" + CordaUtil.getInstance().partyToString(bob.getParty()) + "\",\"issuerRef\":\"100\",\"owner\":\"" + CordaUtil.getInstance().partyToString(charlie.getParty()) + "\"},{\"amt\":{\"currency\":\"USD\", \"quantity\":100}, \"issuer\":\"" + CordaUtil.getInstance().partyToString(alice.getParty()) + "\",\"issuerRef\":\"100\",\"owner\":\"" + CordaUtil.getInstance().partyToString(charlie.getParty()) + "\"}], \"sendPaymentTo\":\"" + CordaUtil.getInstance().partyToString(john.getParty()) + "\", \"sendChangeTo\":\""+ CordaUtil.getInstance().partyToString(charlie.getParty()) + "\", \"paymentAmt\":{\"currency\":\"USD\", \"quantity\":250}}";
			System.out.println(json);
			DocumentContext doc = JsonUtil.getJsonParser().parse(json);
			
			
			context.addInput("input", doc);
			
			context.setContainerService(new CordaContainer(new ArrayList(), "test"));
			
			payment pay = new payment();
			pay.eval(context);
			
			
		}
		
	class MockFlow extends AppFlow {

		public MockFlow(boolean initiating) {
			super(initiating, false);
			// TODO Auto-generated constructor stub
		}

		@Override
		public SignedTransaction call() throws FlowException {
			// TODO Auto-generated method stub
			return null;
		}
		
		
	}

}
