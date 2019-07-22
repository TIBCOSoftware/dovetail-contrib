package com.example.iou;

import java.util.LinkedHashMap;

import com.tibco.dovetail.core.runtime.compilers.AppCompiler;
import com.tibco.dovetail.core.runtime.trigger.ITrigger;

public class IOUContractContractImpl {
	 static LinkedHashMap<String, ITrigger> contractTriggers = AppCompiler.compileApp(IOUContractContractImpl.class.getResourceAsStream("transactions.json"));
    public static ITrigger getTrigger(String name) {
        return contractTriggers.get(name);
    }
}
