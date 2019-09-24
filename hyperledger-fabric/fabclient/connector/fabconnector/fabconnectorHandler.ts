
import {Injectable, Inject, Injector} from "@angular/core";
import {Http} from "@angular/http";
import {Observable} from "rxjs/Observable";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IContributionTypes,
    ActionResult,
    IActionResult,
    WiContribModelService,
    WiContributionUtils,
    IConnectorContribution,
    AUTHENTICATION_TYPE
} from "wi-studio/app/contrib/wi-contrib";
import { IValidationResult, ValidationResult } from "wi-studio/common/models/validation";

@WiContrib({})
@Injectable()
export class fabconnectorHandler extends WiServiceHandlerContribution {

    constructor( @Inject(Injector) injector) {
        super(injector);
    }

    value = (fieldName: string, context: IConnectorContribution): Observable<any> | any => {
        return null;
    }

    validate = (fieldName: string, context: IConnectorContribution): Observable<IValidationResult> | IValidationResult => {
        return ValidationResult.newValidationResult();
    }

    action = (actionId: string, context: IConnectorContribution): Observable<IActionResult> | IActionResult => {
        let aresult: IActionResult = ActionResult.newActionResult();
        let actionResult = {
            context: context,
            authType: AUTHENTICATION_TYPE.BASIC,
            authData: {}
        }
        aresult.setResult(actionResult);
        return aresult
    }
}
