/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.cp;

import net.corda.core.contracts.Amount;
import net.corda.core.contracts.LinearState;
import net.corda.core.contracts.UniqueIdentifier;
import net.corda.core.identity.AbstractParty;
import net.corda.core.identity.Party;
import org.jetbrains.annotations.NotNull;

import com.fasterxml.jackson.annotation.JsonIgnore;

import java.util.Arrays;
import java.util.Currency;
import java.util.List;

public class IOU3 implements LinearState{
    protected Party borrower, lender;
    protected final Amount<Currency> amt;
    protected Amount<Currency> paid;
    protected final UniqueIdentifier linearId;

    public Party getBorrower() {
        return borrower;
    }

    public Party getLender() {
        return lender;
    }

    public void setLender(Party lender) {
        this.lender = lender;
    }

    public Amount<Currency> getAmt() {
        return amt;
    }

    public Amount<Currency> getPaid() {
        return paid;
    }

    public void setPaid(Amount<Currency> paid) {
        this.paid = paid;
    }

    public IOU3(Party lender, Party borrower, Amount<Currency> amt, Amount<Currency> paid, UniqueIdentifier linearId){
        this.borrower = borrower;
        this.lender = lender;
        this.amt = amt;
        this.paid = paid;
        this.linearId = linearId;
    }

    public IOU3(Party lender, Party borrower, Amount<Currency> amt, UniqueIdentifier linearId){
        this(lender, borrower, amt, new Amount<Currency>(0, amt.getToken()), linearId);
    }

    @NotNull
    @Override
    @JsonIgnore
    public List<AbstractParty> getParticipants() {
        return Arrays.asList(borrower, lender);
    }

    @NotNull
    @Override
    public UniqueIdentifier getLinearId() {

        return linearId;
    }

    @Override
    public String toString() {
        return "{" +
                "\"borrower\":" + borrower.getOwningKey().toString() +
                ",\"lender\":" + lender.getOwningKey().toString() +
                ",\"amt\":{\"quantity\":" + amt.getQuantity() + ",\"currency\":\"" + amt.getToken().getCurrencyCode() + "\"}" +
                ",\"paid\":{\"quantity\":" + paid.getQuantity() + ",\"currency\":\"" + paid.getToken().getCurrencyCode() + "\"}" +
                ",\"linearId\":\"" + linearId + "\"" +
                '}';
    }

    @Override
    public boolean equals(Object obj) {
        if(obj instanceof IOU3) {
            IOU3 to = (IOU3) obj;
            return borrower.equals(to.getBorrower()) && lender.equals(to.getLender()) && amt.equals(to.getAmt()) && paid.equals(to.getPaid()) && linearId.equals(to.getLinearId());
        } else {
            return false;
        }
    }
}
