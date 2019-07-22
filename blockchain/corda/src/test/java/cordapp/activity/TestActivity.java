package cordapp.activity;

import static net.corda.finance.Currencies.DOLLARS;
import static org.junit.Assert.*;

import java.io.InputStream;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

import org.junit.Test;

//import com.example.iou.IOU;
import com.fasterxml.jackson.core.JsonProcessingException;

import com.google.common.collect.Lists;
import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.CordaContainer;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.container.cordapp.AppContainer;
import com.tibco.dovetail.container.cordapp.AppFlow;
import com.tibco.dovetail.core.runtime.engine.ContextImpl;
import com.tibco.dovetail.core.runtime.util.JsonUtil;

import cordapp.activity.txnbuilder.txnbuilder;
import cordapp.activity.wallet.wallet;
import junit.framework.Assert;
import net.corda.core.contracts.UniqueIdentifier;
import net.corda.core.flows.FlowException;
import net.corda.core.identity.CordaX500Name;
import net.corda.core.transactions.SignedTransaction;
import net.corda.core.transactions.TransactionBuilder;
import net.corda.testing.core.TestIdentity;
import net.corda.testing.node.MockServices;
import smartcontract_corda.activity.payment.payment;

public class TestActivity {
	
	
	
	@Test
	public void testTxbuilder() throws JsonProcessingException {
		/*
		MockServices mock;
		 TestIdentity bob;
		
		
		bob = new TestIdentity(new CordaX500Name("BigCorp", "New York", "GB"));
		mock = new MockServices(bob);
		
		txnbuilder builder = new txnbuilder();
		
		ContextImpl context = new ContextImpl();
		
		context.addInput("command", "com.example.iou.IssueIOU");
		context.addInput("contractClass", "com.example.iou.IOUContract");
		context.addInput("inputSchema", "[{\"name\":\"iou\",\"type\":\"com.example.iou.IOU\",\"isRef\":false,\"isAsset\":true},{\"name\":\"transactionId\",\"isOptional\":false,\"type\":\"String\"},{\"name\":\"timestamp\",\"isOptional\":false,\"type\":\"DateTime\"}]");
		
		//AppContainer ctnr = new AppContainer(mock);
		AppContainer ctnr = new AppContainer(mock, new MockFlow(true));
		context.setContainerService(ctnr);
		
		String issuer = ctnr.partyToString(bob.getParty());
		DocumentContext doc = JsonUtil.getJsonParser().parse("{}");
	//	DocumentContext iou = JsonUtil.getJsonParser().parse("{\"issuer\":\"" + issuer + "\", \"owner\":\"" + issuer + "\", \"amt\":{\"currency\":\"USD\", \"quantity\":100}, \"linearId\":\"myextid\"}");
		IOU iou = new IOU(bob.getParty(), bob.getParty(), DOLLARS(100), new UniqueIdentifier());
		Map iouvalue = CordaUtil.toJsonObject(iou).json();
		//set linear id to extid only
		iouvalue.put("linearId", "myextid");
		doc.put("$", "iou", iouvalue);
		doc.put("$", "transactionId", "testing");
		doc.put("$", "timestamp", "2019-03-18T12:00:00");
		context.addInput("input", doc);
		builder.eval(context);

		assertEquals(1, ctnr.getFlowService().getOutputStates().size());
		System.out.println("input=" + iou);
		System.out.println("output=" + ctnr.getFlowService().getOutputStates());
		assertEquals(iou.getParticipants(), ctnr.getFlowService().getOutputStates().get(0).getParticipants());
	*/
	}
	
	//@Test
	public void testwallet() {
		MockServices mock;
		TestIdentity bob;
		ContextImpl context = new ContextImpl();
		
		bob = new TestIdentity(new CordaX500Name("BigCorp", "New York", "GB"));
		mock = new MockServices(bob);
		CordaUtil.setServiceHub(mock);
		AppContainer ctnr = new AppContainer(mock, new MockFlow(true));
		context.setContainerService(ctnr);
		
		DocumentContext doc = JsonUtil.getJsonParser().parse("{\"amt\":{\"currency\":\"USD\", quantity:100}, \"issuers\":[\"" + CordaUtil.partyToString(bob.getParty()) + "\"]}");
		
		
		
		context.addInput("operation", "Retrieve Funds");
		context.addInput("input", doc);
		
		wallet wal = new wallet();
		wal.eval(context);
		
		
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
			
			String json = "{\"funds\":[{\"amt\":{\"currency\":\"USD\", quantity:100}, \"issuer\":\"" + CordaUtil.partyToString(bob.getParty()) + "\",\"issuerRef\":\"100\",\"owner\":\"" + CordaUtil.partyToString(charlie.getParty()) + "\"},{\"amt\":{\"currency\":\"USD\", quantity:100}, \"issuer\":\"" + CordaUtil.partyToString(bob.getParty()) + "\",\"issuerRef\":\"100\",\"owner\":\"" + CordaUtil.partyToString(charlie.getParty()) + "\"},{\"amt\":{\"currency\":\"USD\", quantity:100}, \"issuer\":\"" + CordaUtil.partyToString(alice.getParty()) + "\",\"issuerRef\":\"100\",\"owner\":\"" + CordaUtil.partyToString(charlie.getParty()) + "\"}], \"sendPaymentTo\":\"" + CordaUtil.partyToString(john.getParty()) + "\", \"sendChangeTo\":\""+ CordaUtil.partyToString(charlie.getParty()) + "\", \"paymentAmt\":{\"currency\":\"USD\", \"quantity\":250}}";
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
