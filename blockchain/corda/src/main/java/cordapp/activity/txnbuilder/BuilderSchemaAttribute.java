package cordapp.activity.txnbuilder;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import smartcontract.trigger.transaction.model.composer.HLCAttribute;

@JsonIgnoreProperties(ignoreUnknown = true)
public class BuilderSchemaAttribute extends HLCAttribute {
	private boolean isAsset = false;
	private boolean isParty = false;
	
	public boolean isAsset() {
		return isAsset;
	}
	public void setIsAsset(boolean isAsset) {
		this.isAsset = isAsset;
	}
	public boolean isParty() {
		return isParty;
	}
	public void setIsParty(boolean isParty) {
		this.isParty = isParty;
	}
	
}
