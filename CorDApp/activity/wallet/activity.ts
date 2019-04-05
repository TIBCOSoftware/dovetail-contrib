import {Observable} from "rxjs/Observable";
import {Injectable, Injector, Inject} from "@angular/core";
import {Http} from "@angular/http";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    IActivityContribution,
    WiContributionUtils,
    IConnectorContribution
} from "wi-studio/app/contrib/wi-contrib";

import * as lodash from "lodash";

let balschema = "{\"$schema\": \"http://json-schema.org/draft-07/schema#\", \"type\": \"object\", \"required\": [\"currency\"],\"properties\": {\"currency\": {\"type\": \"string\"},\"issuer\": {\"type\": \"array\",\"items\": {\"type\": \"string\"}}}}";
let payschema = "{\"$schema\": \"http://json-schema.org/draft-07/schema#\", \"type\": \"object\", \"required\": [\"payTo\", \"amt\"],\"properties\": {\"payTo\":{\"type\":\"string\"}, \"amt\":{\"type\": \"object\", \"properties\":{\"currency\": {\"type\": \"string\"},\"quantity\":{\"type\":\"number\"}}},\"issuer\": {\"type\": \"array\",\"items\": {\"type\": \"string\"}}}}";                  
let baloutschema = "{\"quantity\":0}";
let payoutschema = "{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"type\":\"array\",\"items\":{\"type\":\"object\",\"properties\":{\"ref\":{\"type\":\"string\"},\"data\":{\"type\":\"object\", \"properties\": {\"amt\":{\"type\":\"object\",\"properties\":{\"currency\":{\"type\":\"string\"},\"quantity\":{\"type\":\"integer\"}}},\"issuer\":{\"type\":\"string\"},\"owner\":{\"type\":\"string\"}}}}}}";
@WiContrib({})
@Injectable()
export class WalletActivityContributionHandler extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http,) {
        super(injector, http);
    }
   
    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {
        let op = context.getField("operation").value;
        switch(fieldName){
            case "input":
                switch(op){
                    case "Account Balance":
                        return balschema;
                    case "Make a Payment":
                        return payschema;
                }
                break;
            case "output":
                switch(op){
                    case "Account Balance":
                        return baloutschema;
                    case "Make a Payment":
                        return payoutschema;
                }
            break;
        }
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        return null;
    }
}