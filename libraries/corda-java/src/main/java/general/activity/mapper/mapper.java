/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package general.activity.mapper;
import com.tibco.dovetail.core.runtime.engine.Context;
import com.tibco.dovetail.core.runtime.activity.IActivity;

public class mapper implements IActivity {

	@Override
	public void eval(Context context) throws IllegalArgumentException {
		Object input = context.getInput("input");
		context.setOutput("output", input);
	}
}
