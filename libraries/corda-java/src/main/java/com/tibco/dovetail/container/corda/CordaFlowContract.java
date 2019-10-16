/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.dovetail.container.corda;

import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.core.runtime.trigger.ITrigger;

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

    protected abstract String getResourceHash();
    protected abstract InputStream getTransactionJson();
	protected abstract ITrigger getTrigger(String name);

    public void verifyTransaction(LedgerTransaction tx) throws IllegalArgumentException {

        Set<PublicKey> allCmdKeys = new HashSet<PublicKey>();
        Set<PublicKey> allStateKeys = new HashSet<PublicKey>();

        tx.getCommands().forEach((CommandWithParties<CommandData> it) -> {
            allCmdKeys.addAll(it.getSigners());
        });

        tx.getInputStates().forEach(it -> {
            it.getParticipants().forEach(p -> {
                allStateKeys.add(p.getOwningKey());
            });
        });

        tx.getOutputs().forEach(it -> {
            it.getData().getParticipants().forEach(p -> {
            		allStateKeys.add(p.getOwningKey());
            });
        });

        ContractsDSL.requireThat(check -> {
            check.using("signatures for all state participants must exist: cmd keys=" + CordaUtil.getInstance().serialize(allCmdKeys) + ", state keys=" + CordaUtil.getInstance().serialize(allStateKeys), allCmdKeys.containsAll(allStateKeys));
            return null;
        });

        tx.getCommands().stream().filter(c -> c.getValue() instanceof CordaCommandDataWithData)
                                 .forEach(c -> {
                                     try {
                                         CordaCommandDataWithData command = (CordaCommandDataWithData)c.getValue();
                                         command.deserialize();
                      
                                         String txName = (String)command.getData("command");
                                         
                                         System.out.println("****** contract " + txName + " verification started ******");
                                         CordaContainer ctnr = new CordaContainer(tx.getInputStates(),  txName);
                                         CordaTransactionService txnSvc = new CordaTransactionService(tx, command);
                                        
                                         getTrigger(txName).invoke(ctnr, txnSvc);

                                         CordaDataService data = (CordaDataService) ctnr.getDataService();
                                         validateOutputs(tx, data.getModifiedStates());
                                         System.out.println("****** contract " + txName + " verified ********");
                                     }catch (Exception e){
                                         throw new IllegalArgumentException(e);
                                     }

                                 });

    }

    private void validateOutputs(LedgerTransaction tx, List<DocumentContext> outputs) throws JsonParseException, JsonMappingException, IOException {
        List<DocumentContext> txOuts = tx.getOutputStates().stream().map(it -> CordaUtil.getInstance().toJsonObject(it)).collect(Collectors.toList());
        CordaUtil.getInstance().compare(txOuts, outputs);
    }
    

    public ContractCommandOutput runCommand(CordaCommandDataWithData command, List<ContractState> inputStates) {
    		try {
    		
             String txName = (String)command.getData("command");
            
             System.out.println("****** run " + txName + " ... ******");
             CordaContainer ctnr = new CordaContainer(inputStates,  txName);
             CordaTransactionService txnSvc = new CordaTransactionService(null, command);
            
             getTrigger(txName).invoke(ctnr, txnSvc);

             CordaDataService data = (CordaDataService) ctnr.getDataService();
             ContractCommandOutput outputs = data.getContractCommandOutput();
            
             System.out.println("****** finish " + txName + ". ********");
             return outputs;
     		
    		}catch(Exception e) {
         	throw new IllegalArgumentException(e);
        }
    }
}
