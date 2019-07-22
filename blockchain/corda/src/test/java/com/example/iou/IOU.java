/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.example.iou;

import net.corda.core.DeleteForDJVM;
import net.corda.core.contracts.BelongsToContract;
import net.corda.core.contracts.ContractState;
import net.corda.core.contracts.SchedulableState;
import net.corda.core.contracts.ScheduledActivity;
import net.corda.core.contracts.StateRef;
import net.corda.core.flows.FlowLogicRefFactory;
import net.corda.core.identity.AbstractParty;
import net.corda.core.serialization.CordaSerializable;

import org.jetbrains.annotations.NotNull;

import com.tibco.dovetail.container.corda.CordaCommandDataWithData;
import com.tibco.dovetail.container.corda.CordaContainer;
import com.tibco.dovetail.container.corda.CordaTransactionService;
import com.tibco.dovetail.core.runtime.flow.ReplyData;
import com.tibco.dovetail.core.runtime.trigger.ITrigger;

import java.util.ArrayList;
import java.util.Currency;
import java.util.List;

@CordaSerializable
@BelongsToContract(com.example.iou.IOUContractContract.class)
public class IOU implements  net.corda.core.contracts.LinearState, SchedulableState {


    private net.corda.core.identity.Party issuer;

    private net.corda.core.identity.Party owner;

    private net.corda.core.contracts.Amount<Currency> amt;

    private net.corda.core.contracts.Amount<Currency> paid;

    private net.corda.core.contracts.UniqueIdentifier linearId;
    
    public IOU () {
    	
    }

    public IOU (net.corda.core.identity.Party issuer,net.corda.core.identity.Party owner,net.corda.core.contracts.Amount<Currency> amt,net.corda.core.contracts.Amount<Currency> paid,net.corda.core.contracts.UniqueIdentifier linearId){

        this.issuer = issuer;

        this.owner = owner;

        this.amt = amt;

        this.paid = paid;

        this.linearId = linearId;

    }


    @DeleteForDJVM 
    @Override
    public ScheduledActivity nextScheduledActivity(StateRef thisStateRef, FlowLogicRefFactory flowLogicRefFactory) {

        ITrigger trig = com.example.iou.IOUContractContractImpl.getTrigger("generateAutoPayEvent");
        if(trig == null)
            return null;
        else {
        		CordaContainer ctnr = new CordaContainer(java.util.Arrays.asList((ContractState)this),  "IOU_nextScheduledActivity");
            ctnr.addContainerProperty("FlowLogicRefFactory", flowLogicRefFactory);
            ctnr.addContainerProperty("StateRef", thisStateRef);
            ctnr.addContainerProperty("Namespace", "com.example.iou");
            
            CordaCommandDataWithData command =  new CordaCommandDataWithData();
            command.putData("transactionInput", this);
            command.putData("command", "generateAutoPayEvent");

            CordaTransactionService txnSvc = new CordaTransactionService(null, command);
                             
            ReplyData reply = trig.invoke(ctnr, txnSvc);
            if(reply == null || reply.getObjectData() == null){
                return null;
            } else {
            		ScheduledActivity sa = (ScheduledActivity)reply.getObjectData();
                return sa;
            }
        }

    }
    public net.corda.core.identity.Party getIssuer() {
        return this.issuer;
    }

    public void setIssuer(net.corda.core.identity.Party issuer) {
        this.issuer = issuer;
    }


    public net.corda.core.identity.Party getOwner() {
        return this.owner;
    }

    public void setOwner(net.corda.core.identity.Party owner) {
        this.owner = owner;
    }


    public net.corda.core.contracts.Amount<Currency> getAmt() {
        return this.amt;
    }

    public void setAmt(net.corda.core.contracts.Amount<Currency> amt) {
        this.amt = amt;
    }


    public net.corda.core.contracts.Amount<Currency> getPaid() {
        return this.paid;
    }

    public void setPaid(net.corda.core.contracts.Amount<Currency> paid) {
        this.paid = paid;
    }



    @NotNull
    @Override
    public net.corda.core.contracts.UniqueIdentifier getLinearId() {
        return this.linearId;
    }

    @NotNull
    @Override
    public List<AbstractParty> getParticipants() {
        List<AbstractParty> participants = new ArrayList<AbstractParty>();
        participants.add(issuer);
        participants.add(owner);
        return participants;
    }

    @Override
    public String toString() {
        return "{\"issuer\":" + issuer.getOwningKey().toString() + ",\"owner\":" + owner.getOwningKey().toString() + ",\"amt\":" + "{\"quantity\":" + amt.getQuantity() + ", \"currency\":\"" + amt.getToken().getCurrencyCode() + "\"}" + ",\"paid\":" + "{\"quantity\":" + paid.getQuantity() + ", \"currency\":\"" + paid.getToken().getCurrencyCode() + "\"}" + ",\"linearId\":" + (linearId.getId().toString() + "_" + linearId.getExternalId())+ "}";
    }

    @Override
    public boolean equals(Object obj) {
        if(obj instanceof IOU) {
            IOU to = (IOU) obj;
            return issuer.equals(to.getIssuer()) && owner.equals(to.getOwner()) && amt.equals(to.getAmt()) && paid.equals(to.getPaid()) && linearId.equals(to.getLinearId());
        } else {
            return false;
        }
    }
}
