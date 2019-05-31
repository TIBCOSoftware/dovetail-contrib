
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
    CreateFlowActionResult,
    WiContributionUtils,
    IActivityContribution,
    IConnectorContribution
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
        if (fieldName === "commonData") {
            return Observable.create(observer => {
                WiContributionUtils.getConnections(this.http, "fabric").subscribe((data: IConnectorContribution[]) => {
                    let connectionRefs = [];
                    data.forEach(connection => {
                        if ((<any>connection).isValid) {
                            for (let i = 0; i < connection.settings.length; i++) {
                                if (connection.settings[i].name === "name") {
                                    connectionRefs.push({
                                        "unique_id": WiContributionUtils.getUniqueId(connection),
                                        "name": connection.settings[i].value
                                    });
                                    break;
                                }
                            }
                        }
                    });
                    observer.next(connectionRefs);
                });
            });
        } else if (fieldName === "parametersDataType" || fieldName === "transientDataType" || fieldName === "returnsDataType") {
            let connectorId = context.getField("commonData").value;
            if (connectorId) {
                return Observable.create(observer => {
                    WiContributionUtils.getConnection(this.http, connectorId).map(data => data)
                    .subscribe(data => {
                        for (let setting of data.settings) {
                            if (setting.name === "schema" && setting.value.value) {
//                                console.log("connector data: " + setting.name + " = " + setting.value.value);
                                let schema = [""];
                                let sobj = JSON.parse(setting.value.value);
                                Object.keys(sobj).forEach( (pkg) => {
                                    Object.keys(sobj[pkg]).forEach( (n) => {
                                        schema.push(pkg + "." + n);
                                    });
                                });
                                observer.next(schema);
                                break;
                            }
                        }
                    });
                });
            }
        } else if (fieldName === "parameters" || fieldName === "transient" || fieldName === "returns") {
            let dataTypeName = context.getField(fieldName+"DataType").value;
            let connectorId = context.getField("commonData").value;
            if (dataTypeName && connectorId) {
                // set pre-defined schema from shared data defs
                return Observable.create(observer => {
                    WiContributionUtils.getConnection(this.http, connectorId).map(data => data)
                    .subscribe(data => {
                        for (let setting of data.settings) {
                            if (setting.name === "schema" && setting.value && setting.value.value) {
//                                console.log(fieldName + " schema data: " + setting.name + " = " + setting.value.value);
                                let sch = JSON.parse(setting.value.value);
                                let idx = dataTypeName.lastIndexOf('.');
                                let pkg = sch[dataTypeName.substring(0, idx)];
                                observer.next(JSON.stringify(pkg[dataTypeName.substring(idx+1)], null, 2));
                                break;
                            }
                        }
                    });
                });
            }
        }
        return null;
    }

    // verify user entries are valid JSON string
    validate = (fieldName: string, context: ITriggerContribution): Observable<IValidationResult> | IValidationResult => {
        let vresult: IValidationResult = ValidationResult.newValidationResult();
        if (fieldName === "commonData") {
            if (context.getMode() === MODE.WIZARD || context.getMode() === MODE.SERVERLESS_FLOW) {
                vresult.setVisible(true);
            } else {
                vresult.setVisible(false);
            }
            return vresult;
        } else if (fieldName === "parametersDataType" || fieldName === "transientDataType" || fieldName === "returnsDataType") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let commonDataField: IFieldDefinition = context.getField("commonData");
            if (commonDataField.value) {
                vresult.setVisible(true);
            } else {
                vresult.setVisible(false);
            }
            return vresult;
        } else if (fieldName === "parameters" || fieldName === "transient" || fieldName === "returns") {
            if (context.getMode() === MODE.WIZARD || context.getMode() === MODE.SERVERLESS_FLOW) {
                let vresult: IValidationResult = ValidationResult.newValidationResult();
                let dataField: IFieldDefinition = context.getField(fieldName);
                let dataTypeField: IFieldDefinition = context.getField(fieldName+"DataType");
                if (dataTypeField.value && dataTypeField.display.visible) {
                    vresult.setReadOnly(true);
                } else {
                    vresult.setReadOnly(false);
                }
                if (dataField.value && dataField.value.value) {
                    try {
                        // verify well-formed JSON schema
                        let valRes;
                        valRes = JSON.parse(dataField.value.value);
                        valRes = JSON.stringify(valRes);
                    } catch (e) {
                        return vresult.setError("FABTIC-TRIGGER-1000", "Invalid JSON: " + e.toString());
                    }
                }
                vresult.setReadOnly(false);
                return vresult;
            } else {
                let vresult: IValidationResult = ValidationResult.newValidationResult();
                vresult.setReadOnly(true);
                return vresult;
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
                                "value": parametersField.value.value,
                                "metadata": ""
                            };
                        } else if (trigger.outputs[j].name === "transient") {
                            trigger.outputs[j].value = {
                                "value": transientField.value.value,
                                "metadata": ""
                            };
                        }
                    }
                }
                if (trigger && trigger.reply && trigger.reply.length > 0) {
                    for (let j = 0; j < trigger.reply.length; j++) {
                        if (trigger.reply[j].name === "returns") {
                            trigger.reply[j].value = {
                                "value": returnsField.value.value,
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
