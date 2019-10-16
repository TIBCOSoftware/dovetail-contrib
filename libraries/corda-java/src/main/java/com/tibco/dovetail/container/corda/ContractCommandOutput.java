package com.tibco.dovetail.container.corda;

import java.security.PublicKey;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import com.jayway.jsonpath.DocumentContext;

import kotlin.Triple;
import net.corda.core.contracts.CommandData;

public class ContractCommandOutput {
	private List<Triple<String, DocumentContext, CommandData>> outputStates = new ArrayList<Triple<String, DocumentContext, CommandData>>() ;
	private Map<CommandData, Set<PublicKey>> embededCommands = new LinkedHashMap<CommandData, Set<PublicKey>>();
	
	public void addOutputState(String state, DocumentContext doc, CommandData cmd) {
		outputStates.add(new Triple<String, DocumentContext, CommandData>(state, doc, cmd));
	}
	
	public void addCommand(CommandData cmd, PublicKey key) {
		Set<PublicKey> keys = embededCommands.get(cmd);
		if(keys == null) {
			keys = new HashSet<PublicKey>();
			embededCommands.put(cmd, keys);
		}
		
		keys .add(key);
	}

	public List<Triple<String, DocumentContext, CommandData>> getOutputStates(){
		return this.outputStates;
	}
	
	public Map<CommandData, Set<PublicKey>> getEmbeddedCommands(){
		return this.embededCommands;
	}
}
