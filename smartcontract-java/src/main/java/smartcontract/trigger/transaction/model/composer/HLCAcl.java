package smartcontract.trigger.transaction.model.composer;

import java.util.List;
import java.util.Map;

public class HLCAcl {
	private List<String> authorizedParties;
	private Map<String, String> conditions;
	
	public List<String> getAuthorizedParties() {
		return authorizedParties;
	}
	public void setAuthorizedParties(List<String> authorizedParties) {
		this.authorizedParties = authorizedParties;
	}
	public Map<String, String> getConditions() {
		return conditions;
	}
	public void setConditions(Map<String, String> conditions) {
		this.conditions = conditions;
	}
	
}
