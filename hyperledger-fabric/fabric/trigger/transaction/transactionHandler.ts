
import {Injectable, Inject, Injector} from "@angular/core";
import {Http} from "@angular/http";
import {Observable} from "rxjs/Observable";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IActionResult,
    ActionResult,
    WiContribModelService,
    ICreateFlowActionContext,
    CreateFlowActionResult
} from "wi-studio/app/contrib/wi-contrib";
import { ITriggerContribution, IFieldDefinition, MODE } from "wi-studio/common/models/contrib";
import { IValidationResult, ValidationResult } from "wi-studio/common/models/validation";
import * as lodash from "lodash";

@WiContrib({})
@Injectable()
export class transactionHandler extends WiServiceHandlerContribution {

    constructor(private injector: Injector, private http: Http, private contribModelService: WiContribModelService) {
        super(injector, http, contribModelService);
    }

    value = (fieldName: string, context: ITriggerContribution): Observable<any> | any => {
        return null;
    }

    // verify user entries are valid JSON string
    validate = (fieldName: string, context: ITriggerContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === "parameters") {
            if (context.getMode() === MODE.WIZARD || context.getMode() === MODE.SERVERLESS_FLOW) {
                let parametersField: IFieldDefinition = context.getField("parameters");
                if (parametersField.value) {
                    try {
                        // verify well-formed JSON schema
                        let valRes;
                        valRes = JSON.parse(parametersField.value);
                        valRes = JSON.stringify(valRes);
                    } catch (e) {
                        return ValidationResult.newValidationResult().setError("FABTIC-TRIGGER-1000", "Invalid JSON: " + e.toString());
                    }
                }
            }
        } else if (fieldName === "transient") {
            if (context.getMode() === MODE.WIZARD || context.getMode() === MODE.SERVERLESS_FLOW) {
                let transientField: IFieldDefinition = context.getField("transient");
                if (transientField.value) {
                    try {
                        // verify well-formed JSON schema
                        let valRes;
                        valRes = JSON.parse(transientField.value);
                        valRes = JSON.stringify(valRes);
                    } catch (e) {
                        return ValidationResult.newValidationResult().setError("FABTIC-TRIGGER-1000", "Invalid JSON: " + e.toString());
                    }
                }
            }
        } else if (fieldName === "returns") {
            if (context.getMode() === MODE.WIZARD || context.getMode() === MODE.SERVERLESS_FLOW) {
                let returnsField: IFieldDefinition = context.getField("returns");
                if (returnsField.value) {
                    try {
                        let valRes;
                        valRes = JSON.parse(returnsField.value);
                        valRes = JSON.stringify(valRes);
                    } catch (e) {
                        return ValidationResult.newValidationResult().setError("FABTIC-TRIGGER-1000", "Invalid JSON: " + e.toString());
                    }
                }
            }
        }
        return null;
    }

    // used to configure trigger with data from "Add a trigger" wizard
    action = (actionId: string, context: ICreateFlowActionContext): Observable<IActionResult> | IActionResult => {
        let modelService = this.getModelService();
        let result = CreateFlowActionResult.newActionResult();
        if (context.handler && context.handler.settings && context.handler.settings.length > 0) {
            let nameField = <IFieldDefinition>context.getField("name");
            let parametersField = <IFieldDefinition>context.getField("parameters");
            let transientField = <IFieldDefinition>context.getField("transient");
            let returnsField = <IFieldDefinition>context.getField("returns");
            if (nameField && nameField.value) {
                let trigger = modelService.createTriggerElement("fabric/fabric-transaction");
                if (trigger && trigger.handler && trigger.handler.settings && trigger.handler.settings.length > 0) {
                    for (let j = 0; j < trigger.handler.settings.length; j++) {
                        if (trigger.handler.settings[j].name === "name") {
                            trigger.handler.settings[j].value = nameField.value;
                        }
                    }
                }
                if (trigger && trigger.outputs && trigger.outputs.length > 0) {
                    for (let j = 0; j < trigger.outputs.length; j++) {
                        if (trigger.outputs[j].name === "parameters") {
                            trigger.outputs[j].value = {
                                "value": parametersField.value,
                                "metadata": ""
                            };
                        } else if (trigger.outputs[j].name === "transient") {
                            trigger.outputs[j].value = {
                                "value": transientField.value,
                                "metadata": ""
                            };
                        }
                    }
                }
                if (trigger && trigger.reply && trigger.reply.length > 0) {
                    for (let j = 0; j < trigger.reply.length; j++) {
                        if (trigger.reply[j].name === "returns") {
                            trigger.reply[j].value = {
                                "value": returnsField.value,
                                "metadata": ""
                            };
                            break;
                        }
                    }
                }
                let flowModel = modelService.createFlow(nameField.value, context.getFlowDescription());
                result = result.addTriggerFlowMapping(lodash.cloneDeep(trigger), lodash.cloneDeep(flowModel));
            }
        }
        let actionResult = ActionResult.newActionResult().setSuccess(true).setResult(result);
        return actionResult;
    }
}
