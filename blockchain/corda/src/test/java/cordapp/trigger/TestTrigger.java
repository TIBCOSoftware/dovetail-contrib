package cordapp.trigger;

import static net.corda.finance.Currencies.DOLLARS;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

import java.io.IOException;
import java.io.InputStream;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;

import org.junit.Test;

import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jayway.jsonpath.DocumentContext;

import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.container.cordapp.AppContainer;
import com.tibco.dovetail.container.cordapp.AppDataService;
import com.tibco.dovetail.container.cordapp.AppFlow;
import com.tibco.dovetail.container.cordapp.AppTransactionService;
import com.tibco.dovetail.core.model.flow.FlowAppConfig;
import com.tibco.dovetail.core.model.flow.HandlerConfig;
import com.tibco.dovetail.core.model.flow.Resources;
import com.tibco.dovetail.core.model.flow.TriggerConfig;
import com.tibco.dovetail.core.runtime.compilers.AppCompiler;
import com.tibco.dovetail.core.runtime.compilers.FlowCompiler;
import com.tibco.dovetail.core.runtime.flow.ActivityTask;
import com.tibco.dovetail.core.runtime.flow.AttributeMapping;
import com.tibco.dovetail.core.runtime.flow.ReplyData;
import com.tibco.dovetail.core.runtime.flow.TransactionFlow;
import com.tibco.dovetail.core.runtime.services.IContainerService;
import com.tibco.dovetail.core.runtime.services.IDataService;
import com.tibco.dovetail.core.runtime.services.IEventService;
import com.tibco.dovetail.core.runtime.services.ILogService;
import com.tibco.dovetail.core.runtime.trigger.ITrigger;

import net.corda.core.contracts.UniqueIdentifier;
import net.corda.core.flows.FlowException;
import net.corda.core.identity.CordaX500Name;
import net.corda.core.transactions.SignedTransaction;
import net.corda.core.transactions.TransactionBuilder;
import net.corda.testing.core.TestIdentity;
import net.corda.testing.node.MockServices;
public class TestTrigger {
	static TestIdentity self ;
	static MockServices mock;
	
	static {
		self = new TestIdentity(new CordaX500Name("self", "New York", "GB"));
		mock = new MockServices(self);
		
	}
	@Test
	public void testInitiator () throws Exception {
		CordaUtil.setServiceHub(mock);
		ObjectMapper mapper = new ObjectMapper();
		InputStream in = this.getClass().getResourceAsStream("iouapp.json");
		
		LinkedHashMap<String, ITrigger> contractTriggers = AppCompiler.compileApp(in).getTriggers();
	 	assertEquals(3, contractTriggers.size());
	 	
	 	assertNotNull(contractTriggers.get("IssueIOUInitiator"));
	 	
	 	ITrigger trigger = contractTriggers.get("IssueIOUInitiator");
	 	
	 	LinkedHashMap<String, Object> args = new LinkedHashMap<String, Object>();
	 	TestIdentity corp = new TestIdentity(new CordaX500Name("BigCorp", "New York", "GB"));
	 	TestIdentity reg = new TestIdentity(new CordaX500Name("Regulator", "New York", "GB"));
	 	
	 	args.put("owner", corp.getParty());
	 	args.put("amt", DOLLARS(100));
	 	args.put("extId", "testing");
	 	args.put("regulator", reg.getParty());
	 	
	 	AppTransactionService txn = new AppTransactionService(args, "IssueIOUInitiator", self.getParty());
	 	Map<String, Object> triggerData = txn.resolveTransactionInput(trigger.getHandler("IssueIOUInitiator").getTxnInputs());
	 	triggerData.forEach((k,v) -> {
	 		if (v instanceof DocumentContext)
	 			System.out.println(k + "=" + ((DocumentContext)v).jsonString());
	 		else
	 			System.out.println(k + "=" + v);
	 	});
	 	
	 	ITrigger s = contractTriggers.get("SettleIOUInitiator");
	 	TransactionFlow f = s.getHandler("SettleIOUInitiator");
	 	ActivityTask t= f.getTask("Mapper");
	 	
	 	t.getInputs().forEach((k,m) -> System.out.println("key=" + k + ", type=" + m.getMappingType()));
	 	List<AttributeMapping> m = (List<AttributeMapping>)t.getInput("input").getMappingValue();
	 	m.forEach(a -> System.out.println("input attr mapping=" + a.getName() + "="+ a.getMappingType()));
	 	
	// 	AppContainer ctnr = new AppContainer(mock, new MockFlow(false));
	// 	contractTriggers.get("IssueIOUInitiator").invoke(ctnr, txn);
	}
	
	//@Test
	public void testschedulable () throws Exception {
		CordaUtil.setServiceHub(mock);
		ObjectMapper mapper = new ObjectMapper();
		InputStream in = this.getClass().getResourceAsStream("iouapp.json");
		
		LinkedHashMap<String, ITrigger> contractTriggers = AppCompiler.compileApp(in).getTriggers();
	 	assertNotNull(contractTriggers.get("autopayment"));
	 	
	 	ITrigger trigger = contractTriggers.get("autopayment");
	 	
	 	LinkedHashMap<String, Object> args = new LinkedHashMap<String, Object>();
	 	TestIdentity issuer = new TestIdentity(new CordaX500Name("Issuer", "New York", "GB"));
	 	TestIdentity owner = new TestIdentity(new CordaX500Name("Owner", "New York", "GB"));
	 	
	 	args.put("transactionInput", new com.example.iou.IOU(issuer.getParty(), owner.getParty(), DOLLARS(100), DOLLARS(0), new UniqueIdentifier()));
	 	
	 	
	 	AppTransactionService txn = new AppTransactionService(args, "autopayment", self.getParty());
	 	Map<String, Object> triggerData = txn.resolveTransactionInput(trigger.getHandler("autopayment").getTxnInputs());
	 	triggerData.forEach((k,v) -> {
	 		if (v instanceof DocumentContext)
	 			System.out.println(k + "=" + ((DocumentContext)v).jsonString());
	 		else
	 			System.out.println(k + "=" + v);
	 	});
	 	
	 	
	 	AppContainer ctnr = new AppContainer(mock, new MockFlow(false));
	 	ReplyData reply = trigger.invoke(ctnr, txn);
	 	System.out.println("rely=" + reply.getData());
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
	
	public class MockContainer implements IContainerService {

		@Override
		public IDataService getDataService() {
			return new AppDataService(mock, new TransactionBuilder(), UUID.randomUUID());
		}

		@Override
		public IEventService getEventService() {
			return new MockEventService();
		}

		@Override
		public ILogService getLogService() {
			return new MockLogService();
		}

	
		public void addContainerProperty(String name, Object v) {
			// TODO Auto-generated method stub
			
		}

		@Override
		public Object getContainerProperty(String name) {
			// TODO Auto-generated method stub
			return null;
		}

		@Override
		public void addContainerAsyncTask(String name, Object v) {
			// TODO Auto-generated method stub
			
		}
		
	}
	
	public class MockDataService implements IDataService {

		@Override
		public DocumentContext putState(String assetName, String assetKey, DocumentContext assetValue) {
			// TODO Auto-generated method stub
			return assetValue;
		}

		@Override
		public DocumentContext getState(String assetName, String assetKey, DocumentContext keyValue) {
			// TODO Auto-generated method stub
			return null;
		}

		@Override
		public DocumentContext deleteState(String assetName, String assetKey, DocumentContext keyValue) {
			// TODO Auto-generated method stub
			return null;
		}

		@Override
		public List<DocumentContext> lookupState(String assetName, String assetKey, DocumentContext keyValue) {
			// TODO Auto-generated method stub
			return null;
		}

		@Override
		public List<DocumentContext> getHistory(String assetName, String assetKey, DocumentContext keyValue) {
			// TODO Auto-generated method stub
			return null;
		}

		@Override
		public List<DocumentContext> queryState(Object query) {
			// TODO Auto-generated method stub
			return null;
		}

		@Override
		public boolean processPayment(DocumentContext assetValue) {
			// TODO Auto-generated method stub
			return true;
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
