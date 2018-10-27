package com.tibco.cp;

import net.corda.core.identity.AbstractParty;
import org.jetbrains.annotations.NotNull;
import java.util.ArrayList;
import java.util.Currency;
import java.util.List;

class IOU2 implements  net.corda.core.contracts.LinearState {


    private net.corda.core.identity.Party lender;

    private net.corda.core.identity.Party borrower;

    private net.corda.core.contracts.Amount<Currency> amt;

    private net.corda.core.contracts.Amount<Currency> paid;

    private net.corda.core.contracts.UniqueIdentifier linearId;

    public IOU2 (net.corda.core.identity.Party lender,net.corda.core.identity.Party borrower,net.corda.core.contracts.Amount<Currency> amt,net.corda.core.contracts.Amount<Currency> paid,net.corda.core.contracts.UniqueIdentifier linearId){

        this.lender = lender;

        this.borrower = borrower;

        this.amt = amt;

        this.paid = paid;

        this.linearId = linearId;

    }


    public net.corda.core.identity.Party getLender() {
        return this.lender;
    }

    public void setLender(net.corda.core.identity.Party lender) {
        this.lender = lender;
    }


    public net.corda.core.identity.Party getBorrower() {
        return this.borrower;
    }

    public void setBorrower(net.corda.core.identity.Party borrower) {
        this.borrower = borrower;
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
        participants.add(lender);
        participants.add(borrower);
        return participants;
    }

    @Override
    public String toString() {
        return "{\"lender\":" + lender.getOwningKey().toString() + ",\"borrower\":" + borrower.getOwningKey().toString() + ",\"amt\":" + "{\"quantity\":" + amt.getQuantity() + ", \"currency\":\"" + amt.getToken().getCurrencyCode() + "\"}" + ",\"paid\":" + "{\"quantity\":" + paid.getQuantity() + ", \"currency\":\"" + paid.getToken().getCurrencyCode() + "\"}" + ",\"linearId\":" + (linearId.getId().toString() + "_" + linearId.getExternalId())+ "}";
    }

    @Override
    public boolean equals(Object obj) {
        if(obj instanceof IOU3) {
            IOU3 to = (IOU3) obj;
            return lender.equals(to.getLender()) && borrower.equals(to.getBorrower()) && amt.equals(to.getAmt()) && paid.equals(to.getPaid()) && linearId.equals(to.getLinearId());
        } else {
            return false;
        }
    }
}
