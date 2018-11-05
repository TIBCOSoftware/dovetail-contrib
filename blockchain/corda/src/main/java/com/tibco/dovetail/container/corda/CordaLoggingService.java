package com.tibco.dovetail.container.corda;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.tibco.dovetail.core.runtime.services.ILogService;

public class CordaLoggingService implements ILogService{
	Logger logger;
	public CordaLoggingService(String name) {
		logger = LoggerFactory.getLogger(name);
	}
	
	@Override
	public void debug(String msg) {
		logger.debug(msg);
	}

	@Override
	public void info(String msg) {
		logger.info(msg);
	}

	@Override
	public void warning(String msg) {
		logger.warn(msg);
	}

	@Override
	public void error(String errCode, String msg, Throwable err) {
		logger.error(errCode + ":" + msg, err);
	}

}
