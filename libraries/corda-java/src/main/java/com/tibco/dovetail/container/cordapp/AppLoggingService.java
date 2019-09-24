/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.dovetail.container.cordapp;

import com.tibco.dovetail.container.corda.CordaLoggingService;
import org.slf4j.Logger;

public class AppLoggingService extends CordaLoggingService{

	public AppLoggingService(String name) {
		super(name);
	}
	
	public AppLoggingService(Logger logger) {
		super(logger);
	}
}
