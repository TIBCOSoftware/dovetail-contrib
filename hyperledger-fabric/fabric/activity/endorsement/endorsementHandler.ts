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
export class endorsementHandler extends WiServiceHandlerContribution {
  constructor(@Inject(Injector) injector) {
    super(injector);
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
    if (fieldName === "role") {
      let vresult: IValidationResult = ValidationResult.newValidationResult();
      let operationField: IFieldDefinition = context.getField("operation");
      let roleField: IFieldDefinition = context.getField("role");
      if (operationField.value && operationField.value === "ADD") {
        if (roleField.display && roleField.display.visible == false) {
          vresult.setVisible(true);
        }
      } else {
        vresult.setVisible(false);
      }
      return vresult;
    } else if (fieldName === "organizations") {
      let vresult: IValidationResult = ValidationResult.newValidationResult();
      let operationField: IFieldDefinition = context.getField("operation");
      let organizationsField: IFieldDefinition = context.getField(
        "organizations"
      );
      if (
        operationField.value &&
        (operationField.value === "ADD" || operationField.value === "DELETE")
      ) {
        if (
          organizationsField.display &&
          organizationsField.display.visible == false
        ) {
          vresult.setVisible(true);
        }
      } else {
        vresult.setVisible(false);
      }
      return vresult;
    } else if (fieldName === "policy") {
      let vresult: IValidationResult = ValidationResult.newValidationResult();
      let operationField: IFieldDefinition = context.getField("operation");
      let policyField: IFieldDefinition = context.getField("policy");
      if (operationField.value && operationField.value === "SET") {
        if (policyField.display && policyField.display.visible == false) {
          vresult.setVisible(true);
        }
      } else {
        vresult.setVisible(false);
      }
      return vresult;
    }
    return null;
  };
}
