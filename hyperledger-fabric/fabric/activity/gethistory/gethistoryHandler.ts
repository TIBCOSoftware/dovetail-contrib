
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
export class gethistoryHandler extends WiServiceHandlerContribution {
    constructor( @Inject(Injector) injector) {
        super(injector);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
        return null;
    }

    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === "result") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let resultField: IFieldDefinition = context.getField("result");
            if (resultField.value && resultField.value.value) {
                try {
                    let valRes;
                    valRes = JSON.parse(resultField.value.value);
                    valRes = JSON.stringify(valRes);
                } catch (e) {
                    vresult.setError("FABTIC-GETHISTORY-1000", "Invalid JSON in output setting: " + e.toString());
                }
            }
            return vresult;
        }
        return null;
    }
}