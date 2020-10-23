import { Injectable, Injector, Inject } from "@angular/core";
import { Http } from "@angular/http";
import { Observable } from "rxjs/Observable";
import {
  WiContrib,
  WiServiceHandlerContribution,
  IContributionTypes,
  ActionResult,
  IActionResult,
  WiContribModelService,
  IFieldDefinition,
  WiContributionUtils,
  IActivityContribution,
  IConnectorContribution,
} from "wi-studio/app/contrib/wi-contrib";
import {
  IValidationResult,
  ValidationResult,
} from "wi-studio/common/models/validation";

@WiContrib({})
@Injectable()
export class putallHandler extends WiServiceHandlerContribution {
  constructor(private injector: Injector, private http: Http) {
    super(injector, http);
  }

  value = (
    fieldName: string,
    context: IActivityContribution
  ): Observable<any> | any => {
    if (fieldName === "data") {
      let dataField: IFieldDefinition = context.getField("data");
      if (dataField.value) {
        let sch = JSON.parse(dataField.value);
        // check schema object type, and convert to key-value type
        if (sch["type"] == "object") {
          let data = {};
          data["type"] = "object";
          data["properties"] = sch["properties"];
          let keyval = {};
          keyval["key"] = { type: "string" };
          keyval["value"] = data;
          let item = {};
          item["type"] = "object";
          item["properties"] = keyval;
          let dataschema = {};
          dataschema["$schema"] = sch["$schema"];
          dataschema["type"] = "array";
          dataschema["items"] = item;
          return JSON.stringify(dataschema, null, 2);
        }
      }
    } else if (fieldName === "result") {
      // set it the same as data if not using shared data defs
      let dataField: IFieldDefinition = context.getField("data");
      if (dataField.value) {
        return dataField.value;
      }
    }
    return null;
  };

  validate = (
    fieldName: string,
    context: IActivityContribution
  ): Observable<IValidationResult> | IValidationResult => {
    if (fieldName === "data") {
      let vresult: IValidationResult = ValidationResult.newValidationResult();
      let dataField: IFieldDefinition = context.getField("data");
      if (dataField.value) {
        try {
          let valRes;
          valRes = JSON.parse(dataField.value);
          valRes = JSON.stringify(valRes);
        } catch (e) {
          vresult.setError(
            "FABRIC-PUTALL-1000",
            "Invalid JSON: " + e.toString()
          );
        }
      } else {
        vresult.setError("FABRIC-PUTALL-1010", "Data schema must not be empty");
      }
      return vresult;
    }
    return null;
  };
}
