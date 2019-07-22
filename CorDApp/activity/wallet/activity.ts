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
let fundschema = "{\"$schema\": \"http://json-schema.org/draft-07/schema#\", \"type\": \"object\", \"required\": [\"amt\"],\"properties\": {\"amt\":{\"type\": \"object\", \"properties\":{\"currency\": {\"type\": \"string\"},\"quantity\":{\"type\":\"number\"}}},\"issuer\": {\"type\": \"array\",\"items\": {\"type\": \"string\"}}}}";                  

let baloutschema = "{\"currency\":\"\", \"quantity\":0}";
let payoutschema = "{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"type\":\"array\",\"items\":{\"type\":\"object\",\"properties\":{\"ref\":{\"type\":\"string\"},\"data\":{\"type\":\"object\", \"properties\": {\"amt\":{\"type\":\"object\",\"properties\":{\"currency\":{\"type\":\"string\"},\"quantity\":{\"type\":\"integer\"}}},\"issuer\":{\"type\":\"string\"},\"issuerRef\":{\"type\":\"string\"},\"owner\":{\"type\":\"string\"}, \"assetId\":{\"type\":\"string\"}}}}}}";
let fundoutschema = "{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"type\":\"array\",\"items\":{\"type\":\"object\",\"properties\":{\"ref\":{\"type\":\"string\"},\"data\":{\"type\":\"object\", \"properties\": {\"amt\":{\"type\":\"object\",\"properties\":{\"currency\":{\"type\":\"string\"},\"quantity\":{\"type\":\"integer\"}}},\"issuer\":{\"type\":\"string\"},\"issuerRef\":{\"type\":\"string\"},\"owner\":{\"type\":\"string\"}, \"assetId\":{\"type\":\"string\"}}}}}}";


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
                if(op == "Account Balance")
                        return balschema;
                else if (op == "Make a Payment")
                        return payschema;
                else if (op == "Retrieve Funds")
                        return fundschema;
                else
                    console.log("error, op is not supported");
                break;
            case "output":
                if(op == "Account Balance")
                    return baloutschema;
                else if (op == "Make a Payment")
                    return payoutschema;
                else if (op == "Retrieve Funds")
                    return fundoutschema;
                else
                    console.log("error, op is not supported");
            break;
        }
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        return null;
    }
}