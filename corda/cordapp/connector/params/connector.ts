import { Inject, Injectable, Injector } from "@angular/core";
import { WiContrib, WiServiceHandlerContribution, WiContributionUtils,AUTHENTICATION_TYPE } from "wi-studio/app/contrib/wi-contrib";
import { IConnectorContribution } from "wi-studio/common/models/contrib";
import { IValidationResult, ValidationResult, ValidationError } from "wi-studio/common/models/validation";
import { IActionResult, ActionResult } from "wi-studio/common/models/contrib";
import { Observable } from "rxjs/Observable";

@Injectable()
@WiContrib({})
export class ParamsConnectorService extends WiServiceHandlerContribution {
    constructor( @Inject(Injector) injector) {
        super(injector);
    }

    value = (fieldName: string, context: IConnectorContribution): Observable<any> | any => {
        return null;    
    }

    validate = (fieldName: string, context: IConnectorContribution): Observable<IValidationResult> | IValidationResult => {
       return null;
    }

    action = (name: string, context: IConnectorContribution): Observable<IActionResult> | IActionResult => {
            return Observable.create(observer => {
                let actionResult = {
                    context: context,
                    authType: AUTHENTICATION_TYPE.BASIC,
                    authData: {}
                }
                observer.next(ActionResult.newActionResult().setResult(actionResult));
            });
        
    }

    getFieldValue(context, fieldName): any {
        for(let i=0; i<context.settings.length; i++){
            if(context.settings[i].name=== fieldName){
                return context.settings[i].value;
            }
        }
    }
}
