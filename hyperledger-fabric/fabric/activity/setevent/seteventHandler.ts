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
  IActivityContribution,
} from "wi-studio/app/contrib/wi-contrib";
import {
  IValidationResult,
  ValidationResult,
} from "wi-studio/common/models/validation";

@WiContrib({})
@Injectable()
export class seteventHandler extends WiServiceHandlerContribution {
  constructor(@Inject(Injector) injector) {
    super(injector);
  }

  value = (
    fieldName: string,
    context: IActivityContribution
  ): Observable<any> | any => {
    if (fieldName === "result") {
      let payloadField: IFieldDefinition = context.getField("payload");
      if (payloadField.value) {
        return payloadField.value;
      }
    }
    return null;
  };

  validate = (
    fieldName: string,
    context: IActivityContribution
  ): Observable<IValidationResult> | IValidationResult => {
    if (fieldName === "payload") {
      let vresult: IValidationResult = ValidationResult.newValidationResult();
      let payloadField: IFieldDefinition = context.getField("payload");
      if (payloadField.value) {
        try {
          let valRes;
          valRes = JSON.parse(payloadField.value);
          valRes = JSON.stringify(valRes);
        } catch (e) {
          vresult.setError(
            "FABRIC-SETEVENT-1010",
            "Invalid JSON in payload: " + e.toString()
          );
        }
      }
      return vresult;
    }
    return null;
  };
}
