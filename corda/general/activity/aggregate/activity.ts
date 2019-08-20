import {Observable} from "rxjs/Observable";
import {Injectable, Injector, Inject} from "@angular/core";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    IActivityContribution
} from "wi-studio/app/contrib/wi-contrib";
import * as lodash from "lodash";

const numericInputSchema  = {
    "type": "array",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "items": {
      "type": "object",
      "properties": {
        "data": {
          "type": "number"
        }
      }
    }
  };

  const doubleInputSchema  = {
    "type": "array",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "items": {
      "type": "object",
      "properties": {
        "data": {
          "type": "string"
        }
      }
    }
  };
const numericOutputSchema = {
    "type": "object",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "required": ["result"],
    "properties": {
      "result": {
        "type": "number"
      }
    }
  };
  const doubleOutputSchema = {
    "type": "object",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "required": ["result"],
    "properties": {
      "result": {
        "type": "string"
      }
    }
  };

@WiContrib({})
@Injectable()
export class AggregateActivityContributionHandler extends WiServiceHandlerContribution {
    constructor( @Inject(Injector) injector) {
        super(injector);
    }

    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {
        let datatype = context.getField("dataType").value;
        let op = context.getField("operation").value;
        let schemaJSON;
        switch(fieldName){
            case "input":
                return Observable.create(observer => {
                    if(datatype == "Double" || op == "AVG")
                        schemaJSON = lodash.cloneDeep(doubleInputSchema);
                    else
                        schemaJSON = lodash.cloneDeep(numericInputSchema);
                    observer.next(JSON.stringify(schemaJSON));
                });
            case "output":
                return Observable.create(observer => {
                    if(datatype == "Double" || op == "AVG")
                        schemaJSON = lodash.cloneDeep(doubleOutputSchema);
                    else
                    schemaJSON = lodash.cloneDeep(numericOutputSchema);
                    observer.next(JSON.stringify(schemaJSON));
                });
            default:
                return null;
        }
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        let action = context.getField("operation").value;
        let datatype = context.getField("dataType").value;
        switch(fieldName){
            case "precision":
            case "scale":
            case "rounding":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    if(action === "avg" || datatype === "Double")
                        vresult.setVisible(true);
                    else
                        vresult.setVisible(false)
                    observer.next(vresult);
                });
            default:
                return null; 
        }
    }
}