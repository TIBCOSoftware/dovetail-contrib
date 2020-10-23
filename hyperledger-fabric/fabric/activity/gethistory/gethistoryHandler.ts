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
export class gethistoryHandler extends WiServiceHandlerContribution {
  constructor(private injector: Injector, private http: Http) {
    super(injector, http);
  }

  value = (
    fieldName: string,
    context: IActivityContribution
  ): Observable<any> | any => {
    if (fieldName === "result") {
      let resultField: IFieldDefinition = context.getField("result");
      if (resultField.value) {
        let sch = JSON.parse(resultField.value);
        // check schema object type, and convert to history type
        if (sch["type"] == "object") {
          let data = {};
          data["type"] = "object";
          data["properties"] = sch["properties"];
          let hist = {};
          hist["txID"] = { type: "string" };
          hist["txTime"] = { type: "string" };
          hist["isDeleted"] = { type: "boolean" };
          hist["value"] = data;
          let item = {};
          item["type"] = "object";
          item["properties"] = hist;
          let histschema = {};
          histschema["$schema"] = sch["$schema"];
          histschema["type"] = "array";
          histschema["items"] = item;
          return JSON.stringify(histschema, null, 2);
        }
      }
    }
    return null;
  };

  validate = (
    fieldName: string,
    context: IActivityContribution
  ): Observable<IValidationResult> | IValidationResult => {
    if (fieldName === "result") {
      let vresult: IValidationResult = ValidationResult.newValidationResult();
      let resultField: IFieldDefinition = context.getField("result");
      if (resultField.value) {
        try {
          let valRes;
          valRes = JSON.parse(resultField.value);
          valRes = JSON.stringify(valRes);
        } catch (e) {
          vresult.setError(
            "FABRIC-GETHISTORY-1000",
            "Invalid JSON in output setting: " + e.toString()
          );
        }
      }
      return vresult;
    }
    return null;
  };
}
