/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.cp;

import java.util.List;

import net.corda.core.contracts.CommandAndState;
import net.corda.core.contracts.OwnableState;
import net.corda.core.identity.AbstractParty;

public class MyOwnableState implements OwnableState {

	@Override
	public List<AbstractParty> getParticipants() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public AbstractParty getOwner() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public CommandAndState withNewOwner(AbstractParty arg0) {
		// TODO Auto-generated method stub
		return null;
	}

}
