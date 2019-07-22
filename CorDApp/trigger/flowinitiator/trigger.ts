/// <amd-dependency path="./common"/>
import {Injectable, Injector} from "@angular/core";
import {Http} from "@angular/http";
import {Observable} from "rxjs/Observable";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    ActionResult,
    IActionResult,
    ICreateFlowActionContext,
    CreateFlowActionResult,
    WiContribModelService,
    WiContributionUtils,
    IConnectorContribution
} from "wi-studio/app/contrib/wi-contrib";
import { ITriggerContribution, IFieldDefinition, IConnectionAllowedValue, MODE } from "wi-studio/common/models/contrib";
import * as lodash from "lodash";
const commonjs = require("./common");

@WiContrib({})
@Injectable()
export class R3FlowInitiatorTriggerHandler extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http, private contribModelService: WiContribModelService) {
        super(injector, http, contribModelService);
    }
    
    value = (fieldName: string, context: ITriggerContribution): Observable<any> | any => {
        switch(fieldName) {
            case "transactionInput":
                let schemaSelection = context.getField("schemaSelection").value;
                if (schemaSelection === "user"){
                    if(Boolean(context.getField("inputParams").value))
                        return commonjs.createFlowInputSchema(context.getField("inputParams").value.value)
                    else
                        return null;
                       
                } else {
                    return Observable.create(observer => {
                        this.getSchemas(schemaSelection).subscribe( schema => {
                            observer.next(schema);
                        });
                    });  
                }
            case "schemaSelection":
                
                let connectionRefs = [];
                connectionRefs.push({
                    "unique_id": "user",
                    "name": "User Defined..."
                });
                return Observable.create(observer => {
                        WiContributionUtils.getConnections(this.http, "CorDApp").subscribe((data: IConnectorContribution[]) => {
                            data.forEach(connection => {
                                if ((<any>connection).isValid) {
                                    for(let i=0; i < connection.settings.length; i++) {
                                        if(connection.settings[i].name === "name"){
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
                
            default:
                return null;
        }
    }

    validate = (fieldName: string, context: ITriggerContribution): Observable<IValidationResult> | IValidationResult => {
        let input = context.getField("hasObservers").value;
        let schemaSelection = context.getField("schemaSelection").value;
        switch (fieldName) {
            case "observerManual":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    vresult.setVisible(input);
                    observer.next(vresult);
                });
            case "inputParams":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    vresult.setVisible(schemaSelection == "user");
                    observer.next(vresult);
                });
        }
        return null;
    }

    action = (actionId: string, context: ICreateFlowActionContext): Observable<IActionResult> | IActionResult => {
       
        let result = CreateFlowActionResult.newActionResult();
        let flows = []
        return Observable.create(observer => {
            this.createFlow(context,  result);                                    
            let actionResult = ActionResult.newActionResult().setSuccess(true).setResult(result);
            observer.next(actionResult);
        });
    }

    createFlow(context, result) : string{
        let modelService = this.getModelService();
        let initrigger = modelService.createTriggerElement("CorDApp/R3FlowInitiator");
        if (initrigger) {
            for (let s = 0; s < initrigger.handler.settings.length; s++) {
                if (initrigger.handler.settings[s].name === "useAnonymousIdentity") {
                    initrigger.handler.settings[s].value = context.getField("useAnonymousIdentity").value;
                } else if (initrigger.handler.settings[s].name === "hasObservers") {
                    initrigger.handler.settings[s].value = context.getField("hasObservers").value;
                } else if (initrigger.handler.settings[s].name === "observerManual") {
                    initrigger.handler.settings[s].value = context.getField("observerManual").value;
                } else if (initrigger.handler.settings[s].name === "observerFlowName") {
                    initrigger.handler.settings[s].value = context.getField("observerFlowName").value;
                } else if (initrigger.handler.settings[s].name === "useExisting") {
                    initrigger.handler.settings[s].value = context.getField("useExisting").value;
                } else if (initrigger.handler.settings[s].name === "schemaSelection") {
                    initrigger.handler.settings[s].value = context.getField("schemaSelection").value;
                } else {
                    let inputp = context.getField("inputParams");
                    initrigger.handler.settings[s].value = {
                        "metadata": "",
                        "value": context.getField("inputParams").value
                    };
                }
            }
            /*
            for (let j = 0; j < initrigger.outputs.length; j++) {
                if (initrigger.outputs[j].name === "transactionInput") {
                    initrigger.outputs[j].value =  {
                        "value": context.getField("transactionInput").value,
                        "metadata": ""
                    };
                    break;
                }
            }*/
        }

        let flowName = context.getFlowName();
        let iniflowModel = modelService.createFlow(flowName, context.getFlowDescription());
        let builder = modelService.createFlowElement("CorDApp/txnbuilder");
        iniflowModel.addFlowElement(builder);
        result = result.addTriggerFlowMapping(lodash.cloneDeep(initrigger), lodash.cloneDeep(iniflowModel));
        return flowName;
    }
    getSchemas(conId):  Observable<any> {
        let schemas = new Map();
        return Observable.create(observer => {
            WiContributionUtils.getConnection(this.http, conId)
                            .map(data => data)
                            .subscribe(data => {
                                let schemas = new Map();
                                for (let setting of data.settings) {
                                    if(setting.name === "inputParams") {
                                        observer.next(commonjs.createFlowInputSchema(setting.value.value));
                                        break;
                                    }
                                }
                            });
                        });
    }
}