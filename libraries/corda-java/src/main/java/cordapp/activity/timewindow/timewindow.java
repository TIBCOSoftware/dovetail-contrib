package cordapp.activity.timewindow;

import java.time.Duration;
import java.time.Instant;
import java.util.LinkedHashMap;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.cordapp.AppFlow;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.engine.Context;

import net.corda.core.contracts.TimeWindow;

public class timewindow implements IActivity {

	@Override
	public void eval(Context context) throws IllegalArgumentException {
		String windowtype = context.getInput("window").toString();
		Object input = context.getInput("input");
		if (input == null)
			throw new IllegalArgumentException("timewindow: input is not mapped");
		
		LinkedHashMap indoc = ((DocumentContext)input).json();
		AppFlow txservice = (AppFlow) context.getContainerService().getContainerProperty("FlowService");
		
		String from, until;
		int duration;
		switch (windowtype) {
		case "Only valid if after...":
			from = indoc.get("from").toString();
			txservice.setTimeWindow(TimeWindow.fromOnly(Instant.parse(from)));
			break;
		case "Only valid if before...":
			until = indoc.get("until").toString();
			txservice.setTimeWindow(TimeWindow.untilOnly(Instant.parse(until)));
			break;
		case "Only valid if between...":
			from = indoc.get("from").toString();
			until = indoc.get("until").toString();
			txservice.setTimeWindow(TimeWindow.between(Instant.parse(from), Instant.parse(until)));
			break;
		case "Only valid for the duration of...":
			duration = Integer.valueOf(indoc.get("durationSeconds").toString());
			Object start = indoc.get("from");
			if(start == null || start.toString().isEmpty())
				txservice.setTimeWindow(TimeWindow.withTolerance(Instant.now(), Duration.ofSeconds(duration)));
			else
				txservice.setTimeWindow(TimeWindow.withTolerance(Instant.parse(start.toString()), Duration.ofSeconds(duration)));
			
			break;
		}
	}

}
