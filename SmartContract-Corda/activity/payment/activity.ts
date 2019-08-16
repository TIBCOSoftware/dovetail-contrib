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

let input1schema = "{\"$schema\": \"http://json-schema.org/draft-07/schema#\",\"type\": \"array\",\"items\": {\"type\": \"object\",\"required\": [],\"properties\": {\"issuer\": {\"type\": \"string\"}}}}"
let input2schema = "{\"$schema\": \"http://json-schema.org/draft-07/schema#\",\"type\": \"object\",\"required\": [],\"properties\": {\"cashIssuer\": {\"type\": \"party\"}, \"paidTo\": {\"type\": \"party\"}}}"

@WiContrib({})
@Injectable()
export class PaymentActivityContributionHandler extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http,) {
        super(injector, http);
    }
   
    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        return null;
    }
}