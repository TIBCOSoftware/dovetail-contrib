
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
    WiContributionUtils,
    IActivityContribution,
    IConnectorContribution
} from "wi-studio/app/contrib/wi-contrib";
import { IValidationResult, ValidationResult } from "wi-studio/common/models/validation";

@WiContrib({})
@Injectable()
export class cidHandler extends WiServiceHandlerContribution {
    constructor( private injector: Injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
        return null;
    }

    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === "attrs") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let attrsField: IFieldDefinition = context.getField("attrs");
            if (attrsField.value && attrsField.value.value) {
                try {
                    let valRes;
                    valRes = JSON.parse(attrsField.value.value);
                    valRes = JSON.stringify(valRes);
                } catch (e) {
                    vresult.setError("FABTIC-CID-1000", "Invalid JSON: " + e.toString());
                }
            }
            return vresult;
        }
        return null;
    }
}
