package cordapp.trigger;

import static net.corda.finance.Currencies.DOLLARS;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

import java.io.IOException;
import java.io.InputStream;
import java.util.LinkedHashMap;
import java.util.Map;

import org.junit.Test;

import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.container.cordapp.AppContainer;
import com.tibco.dovetail.container.cordapp.AppFlow;
import com.tibco.dovetail.container.cordapp.AppTransactionService;
import com.tibco.dovetail.core.model.flow.FlowAppConfig;
import com.tibco.dovetail.core.model.flow.HandlerConfig;
import com.tibco.dovetail.core.model.flow.Resources;
import com.tibco.dovetail.core.model.flow.TriggerConfig;
import com.tibco.dovetail.core.runtime.compilers.FlowCompiler;
import com.tibco.dovetail.core.runtime.engine.DovetailEngine;
import com.tibco.dovetail.core.runtime.flow.TransactionFlow;
import com.tibco.dovetail.core.runtime.trigger.ITrigger;

import net.corda.core.flows.FlowException;
import net.corda.core.identity.CordaX500Name;
import net.corda.core.transactions.SignedTransaction;
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
		
		FlowAppConfig app = FlowAppConfig.parseModel(in);
	 	DovetailEngine engine = new DovetailEngine(app);
	 	LinkedHashMap<String, ITrigger> contractTriggers = engine.getTriggers();
	 	assertEquals(2, contractTriggers.size());
	 	
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
	 	Map<String, Object> triggerData = txn.resolveTransactionInput(trigger.getHandler("IssueIOUInitiator").getFlowInputs());
	 	triggerData.forEach((k,v) -> {
	 		if (v instanceof DocumentContext)
	 			System.out.println(k + "=" + ((DocumentContext)v).jsonString());
	 		else
	 			System.out.println(k + "=" + v);
	 	});
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
