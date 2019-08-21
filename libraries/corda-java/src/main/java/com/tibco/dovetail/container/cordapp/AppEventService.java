/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.dovetail.container.cordapp;

import com.tibco.dovetail.core.runtime.services.IEventService;

public class AppEventService implements IEventService {

	@Override
	public void publish(String evtName, String metadata, String evtPayload) {
		return;
	}

}
