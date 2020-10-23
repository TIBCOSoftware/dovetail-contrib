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
export class getHandler extends WiServiceHandlerContribution {
  constructor(private injector: Injector, private http: Http) {
    super(injector, http);
  }

  value = (
    fieldName: string,
    context: IActivityContribution
  ): Observable<any> | any => {
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
          vresult.setError("FABRIC-GET-1000", "Invalid JSON: " + e.toString());
        }
      }
      return vresult;
    }
    return null;
  };
}
