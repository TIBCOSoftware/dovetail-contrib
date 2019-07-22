package cordapp.activity.wallet;

import java.util.Currency;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Set;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.container.cordapp.AppDataService;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.engine.Context;

import net.corda.core.contracts.Amount;
import net.corda.core.contracts.StateAndRef;
import net.corda.core.identity.AbstractParty;
import net.corda.finance.contracts.asset.Cash;

public class wallet implements IActivity{

	@Override
	public void eval(Context ctx) throws IllegalArgumentException {
		String op = ctx.getInput("operation").toString();
		
		LinkedHashMap<String, Object> val = getInputValue(ctx);	
		Object issuers = val.get("issuers");
		
		AppDataService dataservice = (AppDataService) ctx.getContainerService().getDataService();
		
		switch (op) {
			case "Account Balance":
				Object currency = val.get("currency");
				Amount<Currency> bal = dataservice.getAccountBalance(Currency.getInstance(currency.toString()));
				ctx.setOutput("output", CordaUtil.toJsonObject(bal));
				break;
			case "Retrieve Funds":
				Set<AbstractParty> pIssuers = new HashSet<AbstractParty>();
				if(issuers != null) {
					((List)issuers).forEach(i -> pIssuers.add(CordaUtil.partyFromString(i.toString())));
				}
				
				List<StateAndRef<Cash.State>> funds = dataservice.getFunds(pIssuers, getRetriveAmt(val));
				ctx.setOutput("output", CordaUtil.toJsonObject(funds));
		}
	}
	
	private LinkedHashMap<String, Object> getInputValue(Context ctx) {
		Object input = ctx.getInput("input");
		if(input == null) {
			throw new IllegalArgumentException("wallet:: input must be set");
		}
		
		return ((DocumentContext) input).json();
	}
	
	private Amount<Currency> getRetriveAmt(LinkedHashMap<String, Object> val) {
		Object amt = val.get("amt");
		if (amt == null)
			throw new IllegalArgumentException("wallet::Retreive Funds - must specifiy an amount to get");
		
		LinkedHashMap<String, Object> amtval = (LinkedHashMap<String, Object>)amt;
		return new Amount<Currency>(Long.valueOf(amtval.get("quantity").toString()), Currency.getInstance(amtval.get("currency").toString()));
	}

}
