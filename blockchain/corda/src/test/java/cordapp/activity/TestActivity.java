package cordapp.activity;

import static net.corda.finance.Currencies.DOLLARS;
import static org.junit.Assert.*;

import java.io.InputStream;

import org.junit.Test;

import com.example.iou.IOU;
import com.fasterxml.jackson.core.JsonProcessingException;

import com.google.common.collect.Lists;
import com.jayway.jsonpath.DocumentContext;

import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.container.cordapp.AppContainer;
import com.tibco.dovetail.container.cordapp.AppFlow;
import com.tibco.dovetail.core.runtime.engine.ContextImpl;
import com.tibco.dovetail.core.runtime.util.JsonUtil;

import cordapp.activity.txnbuilder.txnbuilder;
import junit.framework.Assert;
import net.corda.core.contracts.UniqueIdentifier;
import net.corda.core.flows.FlowException;
import net.corda.core.identity.CordaX500Name;
import net.corda.core.transactions.SignedTransaction;
import net.corda.core.transactions.TransactionBuilder;
import net.corda.testing.core.TestIdentity;
import net.corda.testing.node.MockServices;

public class TestActivity {
	
	
	
	//@Test
	public void testTxbuilder() throws JsonProcessingException {
		MockServices mock;
		 TestIdentity bob;
		
		
		bob = new TestIdentity(new CordaX500Name("BigCorp", "New York", "GB"));
		mock = new MockServices(bob);
		
		txnbuilder builder = new txnbuilder();
		
		ContextImpl context = new ContextImpl();
		
		context.addInput("command", "com.example.iou.IssueIOU");
		context.addInput("contractClass", "com.example.iou.IOUContract");
		
		AppContainer ctnr = new AppContainer(mock);
		context.setContainerService(ctnr);
		
		DocumentContext doc = JsonUtil.getJsonParser().parse("{}");
		IOU iou = new IOU(bob.getParty(), bob.getParty(), DOLLARS(100), new UniqueIdentifier());
		doc.put("$", "iou", CordaUtil.toJsonObject(iou).json());
		doc.put("$", "transactionId", "testing");
		doc.put("$", "timestamp", "2019-03-18T12:00:00");
		context.addInput("input", doc);
		builder.eval(context);

		assertEquals(1, ctnr.getFlowService().getOutputStates().size());
		System.out.println("input=" + iou);
		System.out.println("output=" + ctnr.getFlowService().getOutputStates());
		assertEquals(iou, ctnr.getFlowService().getOutputStates().get(0));
	}
	
	class MockFlow extends AppFlow {

		public MockFlow(boolean initiating) {
			super(initiating);
			// TODO Auto-generated constructor stub
		}

		@Override
		public SignedTransaction call() throws FlowException {
			// TODO Auto-generated method stub
			return null;
		}
		
	}

}
