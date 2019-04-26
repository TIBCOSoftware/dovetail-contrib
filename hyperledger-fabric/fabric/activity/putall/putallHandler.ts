
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
export class putallHandler extends WiServiceHandlerContribution {
    constructor( @Inject(Injector) injector) {
        super(injector);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
        if (fieldName === "result") {
            let dataField: IFieldDefinition = context.getField("data");
            if (dataField && dataField.value && dataField.value.value) {
                return dataField.value.value;
            }
        }
        return null;
    }

    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === "collection") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let isPrivateField: IFieldDefinition = context.getField("isPrivate");
            let collectionField: IFieldDefinition = context.getField("collection");
            if (isPrivateField.value && isPrivateField.value === true) {
                if (collectionField.display && collectionField.display.visible == false) {
                    vresult.setVisible(true);
                }
            } else {
                vresult.setVisible(false);
            }
            return vresult;
        } else if (fieldName === "data") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let dataField: IFieldDefinition = context.getField("data");
            if (dataField.value && dataField.value.value) {
                try {
                    let valRes;
                    valRes = JSON.parse(dataField.value.value);
                    valRes = JSON.stringify(valRes);
                } catch (e) {
                    vresult.setError("FABTIC-PUTALL-1000", "Invalid JSON: " + e.toString());
                }
            } else {
                vresult.setError("FABTIC-PUTALL-1010", "Data schema must not be empty");
            }
            return vresult;
        }
        return null;
    }
}
