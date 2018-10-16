import { Http } from "@angular/http";
import { Injectable, Inject, Injector } from "@angular/core";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IActionResult,
    IContributionTypes,
    WiContribModelService,
    ICreateFlowActionContext,
    CreateFlowActionResult,
    ActionResult,
    IApplication,
    WiContributionUtils,
    APP_DEPLOYMENT
} from "wi-studio/app/contrib/wi-contrib";
import { ITriggerContribution, IFieldDefinition, MODE } from "wi-studio/common/models/contrib";
import { IValidationResult, ValidationResult } from "wi-studio/common/models/validation";
import { Observable } from "rxjs/Observable";
import * as lodash from "lodash";

@Injectable()
@WiContrib({})

export class JsonSmartContractTriggerService extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http, private contribModelService: WiContribModelService) {
        super(injector, http, contribModelService);
    }

    value = (fieldName: string, context: ITriggerContribution): Observable<any> | any => {
        return null;
    }

    validate = (fieldName: string, context: ITriggerContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === "body") {
            if (context.getMode() === MODE.WIZARD) {
                let body: IFieldDefinition = context.getField("body");
                let valRes;
                if (body.value) {
                    try {
                        valRes = JSON.parse(body.value);
                        valRes = JSON.stringify(valRes);
                    } catch (e) {
                        return ValidationResult.newValidationResult().setError("SCHEMA_ERROR", "Unexpected string in JSON");
                    }
                }
            } else {
                return ValidationResult.newValidationResult().setVisible(true);
            }
        }
        return null;
    }

    action = (fieldName: string, context: ICreateFlowActionContext): Observable<IActionResult> | IActionResult => {
        let modelService = this.getModelService();
        let result = CreateFlowActionResult.newActionResult();
        if (context.handler && context.handler.settings && context.handler.settings.length > 0) {
            let output = <IFieldDefinition>context.getField("body");
            let trigger = modelService.createTriggerElement("SmartContract/tibco-json-smartcontract");
            if (trigger && trigger.outputs && trigger.outputs.length > 0) {
                for (let j = 0; j < trigger.outputs.length; j++) {
                    if (trigger.outputs[j].name === "body") {
                        trigger.outputs[j].value = {
                            "value": output.value,
                            "metadata": ""
                        };
                        break;
                    }
                }
            }
            let flowModel = modelService.createFlow(context.getFlowName(), context.getFlowDescription());
            result = result.addTriggerFlowMapping(lodash.cloneDeep(trigger), lodash.cloneDeep(flowModel));
        }
        let actionResult = ActionResult.newActionResult().setSuccess(true).setResult(result);
        return actionResult;
    }

    formSettings(applicationModel: IApplication): any[] {
        let settings: any[] = [];
        return settings;
    }
}
