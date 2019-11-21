package com.example.cp;

import com.tibco.dovetail.core.runtime.compilers.AppCompiler;
import com.tibco.dovetail.core.runtime.compilers.App;
import com.tibco.dovetail.core.runtime.trigger.ITrigger;

public class IOUContractImpl {
	static App app = AppCompiler.compileApp(IOUContractImpl.class.getResourceAsStream("transactions.json"));
    public static ITrigger getTrigger(String name) {
        return app.getTriggers().get(name);
    }
}
