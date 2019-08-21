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

let inputschema = "{\"$schema\": \"http://json-schema.org/draft-07/schema#\",\"type\": \"object\",\"properties\": {\"txnHash\": {\"type\": \"string\"}, \"index\":{\"type\":\"integer\"}}}"

@WiContrib({})
@Injectable()
export class TimeWindowActivityContributionHandler extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http,) {
        super(injector, http);
    }
   
    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {
        let window = context.getField("window").value;
       
        switch(fieldName) {
            case "input":
                if(Boolean(window) == false)
                    return null;

                switch(window){
                    case "Only valid if after...":
                        return "{\"$schema\": \"http://json-schema.org/draft-07/schema#\",\"type\": \"object\",\"properties\": {\"from\": {\"type\": \"string\", \"format\":\"date-time\"}}, \"required\":[\"from\"]}"
                    case "Only valid if before...":
                        return "{\"$schema\": \"http://json-schema.org/draft-07/schema#\",\"type\": \"object\",\"properties\": {\"until\": {\"type\": \"string\", \"format\":\"date-time\"}}, \"required\":[\"until\"]}"
                    case "Only valid if between...":
                        return "{\"$schema\": \"http://json-schema.org/draft-07/schema#\",\"type\": \"object\",\"properties\": {\"from\": {\"type\": \"string\", \"format\":\"date-time\"}, \"until\": {\"type\": \"string\", \"format\":\"date-time\"}}, \"required\":[\"from\", \"until\"]}"
                    case "Only valid for the duration of...":
                        return "{\"$schema\": \"http://json-schema.org/draft-07/schema#\",\"type\": \"object\",\"properties\": {\"from\": {\"type\": \"string\", \"format\":\"date-time\"},\"durationSeconds\": {\"type\": \"integer\"}}, \"required\":[\"durationSeconds\"]}"
                }
        }
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        return null;
    }

    
}