package com.tibco.dovetail.container.cordapp;

import net.corda.core.node.services.vault.QueryCriteria.CommonQueryCriteria;

public class VaultQuery {
	private String state;
	private CommonQueryCriteria criteria;
	
	public VaultQuery(String state, CommonQueryCriteria criteria) {
		this.state = state;
		this.criteria = criteria;
	}

	public String getState() {
		return state;
	}

	public CommonQueryCriteria getCriteria() {
		return criteria;
	}
}
