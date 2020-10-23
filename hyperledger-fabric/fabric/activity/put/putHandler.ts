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
export class putHandler extends WiServiceHandlerContribution {
  constructor(private injector: Injector, private http: Http) {
    super(injector, http);
  }

  value = (
    fieldName: string,
    context: IActivityContribution
  ): Observable<any> | any => {
    if (fieldName === "result") {
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
        let valRes;
        try {
          valRes = JSON.parse(dataField.value);
          valRes = JSON.stringify(valRes);
        } catch (e) {
          vresult.setError("FABRIC-PUT-1020", "Invalid JSON: " + e.toString());
        }
      } else {
        vresult.setError(
          "FABRIC-PUT-1010",
          "Data definition must not be empty"
        );
      }
      return vresult;
    }
    return null;
  };
}
