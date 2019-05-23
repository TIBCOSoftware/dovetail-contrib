/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.dovetail.container.corda;

import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.core.model.flow.FlowAppConfig;
import com.tibco.dovetail.core.runtime.engine.DovetailEngine;
import com.tibco.dovetail.core.runtime.trigger.ITrigger;

import kotlin.Triple;
import net.corda.core.contracts.*;
import net.corda.core.serialization.CordaSerializable;
import net.corda.core.transactions.LedgerTransaction;

import java.io.IOException;
import java.io.InputStream;
import java.security.PublicKey;
import java.util.*;
import java.util.stream.Collectors;

@CordaSerializable
public abstract class CordaFlowContract {
    private static LinkedHashMap<String, ITrigger> contractTriggers = null;

    protected abstract String getResourceHash();
    protected abstract InputStream getTransactionJson();

    public void verifyTransaction(LedgerTransaction tx) throws IllegalArgumentException {

        Set<PublicKey> allCmdKeys = new HashSet<PublicKey>();
        Set<PublicKey> allStateKeys = new HashSet<PublicKey>();

        tx.getCommands().forEach((CommandWithParties<CommandData> it) -> {
            allCmdKeys.addAll(it.getSigners());
        });

        tx.getInputStates().forEach(it -> {
            it.getParticipants().forEach(p -> {
               // if(!(p instanceof AnonymousParty))
                    allStateKeys.add(p.getOwningKey());
            });
        });

        tx.getOutputs().forEach(it -> {
            it.getData().getParticipants().forEach(p -> {
              //  if(!(p instanceof AnonymousParty))
                    allStateKeys.add(p.getOwningKey());
            });
        });

        ContractsDSL.requireThat(check -> {
            check.using("signatures for all state participants must exist: cmd keys=" + allCmdKeys.toString() + ", state keys=" + allStateKeys.toString(), allCmdKeys.containsAll(allStateKeys));
            return null;
        });

       compileAndCacheTrigger();

        tx.getCommands().stream().filter(c -> c.getValue() instanceof CordaCommandDataWithData)
                                 .forEach(c -> {
                                     try {
                                         CordaCommandDataWithData command = (CordaCommandDataWithData)c.getValue();
                                         command.deserialize();
                      
                                         String txName = (String)command.getData("command");
                                         
                                         System.out.println("****** contract " + txName + " verification started ******");
                                         CordaContainer ctnr = new CordaContainer(tx.getInputStates(),  txName);
                                         CordaTransactionService txnSvc = new CordaTransactionService(tx, command);
                                        
                                         contractTriggers.get(txName).invoke(ctnr, txnSvc);

                                         CordaDataService data = (CordaDataService) ctnr.getDataService();
                                         validateOutputs(tx, data.getModifiedStates());
                                         System.out.println("****** contract " + txName + " verified ********");
                                     }catch (Exception e){
                                         throw new IllegalArgumentException(e);
                                     }

                                 });

    }

    private void validateOutputs(LedgerTransaction tx, List<DocumentContext> outputs) throws JsonParseException, JsonMappingException, IOException {
        List<DocumentContext> txOuts = tx.getOutputStates().stream().map(it -> CordaUtil.toJsonObject(it)).collect(Collectors.toList());
        CordaUtil.compare(txOuts, outputs);
    }
    
    private synchronized void compileAndCacheTrigger() {
    	 try {
 	        //compile flow app and cache the trigger object
    		   InputStream txJson = getTransactionJson();
 	        if(contractTriggers == null) {
 	        	 	
 	        	 	FlowAppConfig app = FlowAppConfig.parseModel(txJson);
 	        	 	DovetailEngine engine = new DovetailEngine(app);
 	        	 	contractTriggers = engine.getTriggers();
 	        }
         }catch(Exception e) {
         		throw new IllegalArgumentException(e);
         }
    }
    
    public List<Triple<String, DocumentContext, CommandData>> runCommand(CordaCommandDataWithData command, List<ContractState> inputStates) {
    		try {
    			 compileAndCacheTrigger();
    	
             String txName = (String)command.getData("command");
            
             System.out.println("****** run " + txName + " ... ******");
             CordaContainer ctnr = new CordaContainer(inputStates,  txName);
             CordaTransactionService txnSvc = new CordaTransactionService(null, command);
            
             contractTriggers.get(txName).invoke(ctnr, txnSvc);

             CordaDataService data = (CordaDataService) ctnr.getDataService();
             List<Triple<String,DocumentContext, CommandData>> outputs = data.getModifiedStatesAndNames();
            
             System.out.println("****** finish " + txName + ". ********");
             return outputs;
     		
    		}catch(Exception e) {
         	throw new IllegalArgumentException(e);
        }
    }
}
