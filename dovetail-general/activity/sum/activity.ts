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

const GenericSchema = {
  "type": "object",
  "$schema": "http://json-schema.org/draft-07/schema#",
  "properties":{
      "dataset":{
          "type":"array",
          "items": {
            "type":"object",
            "properties":{
              "field":{"type":"number"}
            }
          }
      }
  },
  "required":["dataset"]
};

const OutputSchema = {
    "type": "object",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
      "result": {
        "type": "number"
      }
    }
  };

  const GroupbyOutputSchema = {
    "type": "object",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "properties": {
      "result": {
        "type": "array",
        "items":{
          "type":"object",
          "properties":{
            "groupBy":{"type":"string"},
            "value":{"type":"number"}
          }
        }
      }
    }
  };
  

@WiContrib({})
@Injectable()
export class SumActivityContributionHandler extends WiServiceHandlerContribution {
    constructor( @Inject(Injector) injector) {
        super(injector);
    }

    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {
        let groupBy = context.getField("groupby").value;
        let schemaJSON;
        switch(fieldName){
            case "input":
                return Observable.create(observer => {
                    schemaJSON = lodash.cloneDeep(GenericSchema);
                    if(groupBy){
                      schemaJSON.properties["groupBy"] = {type:"string"}
                      schemaJSON.required.push("groupBy")
                    }

                    observer.next(JSON.stringify(schemaJSON));
                });
            case "output":
                return Observable.create(observer => {
                    if(groupBy){
                       schemaJSON = lodash.cloneDeep(GroupbyOutputSchema);
                    } else {
                      schemaJSON = lodash.cloneDeep(OutputSchema);
                    }
                    observer.next(JSON.stringify(schemaJSON));
                });
            default:
                return null;
        }
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        return null; 
    }
}