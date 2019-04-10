package cordapp.activity.wallet;

import java.util.Currency;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.container.cordapp.AppContainer;
import com.tibco.dovetail.container.cordapp.AppDataService;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.engine.Context;
import com.tibco.dovetail.core.runtime.util.JsonUtil;

import net.corda.core.contracts.Amount;
import net.corda.core.contracts.StateAndRef;
import net.corda.core.identity.AbstractParty;
import net.corda.finance.contracts.asset.Cash;

public class wallet implements IActivity{

	@Override
	public void eval(Context ctx) throws IllegalArgumentException {
		String op = ctx.getInput("operation").toString();
		Object input = ctx.getInput("input");
		Object issuers = null;
		LinkedHashMap val = null;
		
		 if(input == null) {
			throw new RuntimeException("wallet:: input must be set");
		 }
		 
		 val = ((DocumentContext) input).json();
		 issuers = val.get("issuers");
		AppDataService dataservice = (AppDataService) ctx.getContainerService().getDataService();
		AppContainer container = (AppContainer) ctx.getContainerService();
		
		switch (op) {
			case "Account Balance":
				Object currency = val.get("currency");
				Amount<Currency> bal = dataservice.getAccountBalance(Currency.getInstance(currency.toString()));
				ctx.setOutput("output", CordaUtil.toJsonObject(bal));
				break;
			case "Retrieve Funds":
				Set<AbstractParty> pIssuers = new HashSet<AbstractParty>();
				if(issuers != null) {
					((List)issuers).forEach(i -> pIssuers.add(container.partyFromString(i.toString())));
				}
				
				Object amt = val.get("amt");
				if (amt == null)
					throw new RuntimeException("wallet::Retreive Funds - must specifiy an amount to get");
				
				LinkedHashMap amtval = (LinkedHashMap)amt;
				List<StateAndRef<Cash.State>> funds = dataservice.getFunds(pIssuers,  new Amount<Currency>(Long.valueOf(amtval.get("quantity").toString()), Currency.getInstance(amtval.get("currency").toString())));
				ctx.setOutput("output", CordaUtil.toJsonObject(funds));
		}
	}

}
