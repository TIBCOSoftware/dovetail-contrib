
import {Injectable, Injector, Inject} from "@angular/core";
import {Http} from "@angular/http";
import {Observable} from "rxjs/Observable";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IContributionTypes,
    ActionResult,
    IActionResult,
    WiContribModelService,
    IFieldDefinition,
    IActivityContribution
} from "wi-studio/app/contrib/wi-contrib";
import { IValidationResult, ValidationResult } from "wi-studio/common/models/validation";

@WiContrib({})
@Injectable()
export class invokechaincodeHandler extends WiServiceHandlerContribution {
    constructor( @Inject(Injector) injector) {
        super(injector);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
        return null;
    }

    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === "parameters") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let parametersField: IFieldDefinition = context.getField("parameters");
            if (parametersField.value && parametersField.value.value) {
                try {
                    let valRes;
                    valRes = JSON.parse(parametersField.value.value);
                    valRes = JSON.stringify(valRes);
                } catch (e) {
                    vresult.setError("FABTIC-INVOKE-1000", "Invalid JSON: " + e.toString());
                }
            }
            return vresult;
        } else if (fieldName === "result") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let resultField: IFieldDefinition = context.getField("result");
            if (resultField.value && resultField.value.value) {
                try {
                    let valRes;
                    valRes = JSON.parse(resultField.value.value);
                    valRes = JSON.stringify(valRes);
                } catch (e) {
                    vresult.setError("FABTIC-INVOKE-1010", "Invalid JSON: " + e.toString());
                }
            }
            return vresult;
        }
        return null;
    }
}