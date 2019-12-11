package com.example.cp;

import net.corda.core.identity.AbstractParty;

import java.util.ArrayList;
import java.util.Currency;
import java.util.List;

import org.jetbrains.annotations.NotNull;

import net.corda.core.serialization.CordaSerializable;
import com.fasterxml.jackson.annotation.JsonIgnore;
import net.corda.core.contracts.*;

import com.tibco.dovetail.container.corda.*;




@CordaSerializable
@BelongsToContract(IOUContract.class)

public class IOU implements LinearState {
   
	private net.corda.core.identity.Party issuer;

    private net.corda.core.identity.Party owner;

    private net.corda.core.contracts.Amount<Currency> amt;

    private net.corda.core.contracts.Amount<Currency> paid;

    private net.corda.core.contracts.UniqueIdentifier linearId;
    private java.time.Instant createDt;
    
    public IOU () {
    	
    }

    public IOU (java.time.Instant dt, net.corda.core.identity.Party issuer,net.corda.core.identity.Party owner,net.corda.core.contracts.Amount<Currency> amt,net.corda.core.contracts.Amount<Currency> paid,net.corda.core.contracts.UniqueIdentifier linearId){

        this.issuer = issuer;

        this.owner = owner;

        this.amt = amt;

        this.paid = paid;

        this.linearId = linearId;
        this.createDt = dt;

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



    public java.time.Instant getCreateDt() {
		return createDt;
	}

	public void setDt(java.time.Instant dt) {
		this.createDt = dt;
	}

	@NotNull
    @Override
    public net.corda.core.contracts.UniqueIdentifier getLinearId() {
        return this.linearId;
    }

    @NotNull
    @Override
    @JsonIgnore
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
