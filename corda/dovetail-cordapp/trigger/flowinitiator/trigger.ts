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
                
                if(Boolean(context.getField("inputParams").value))
                    return commonjs.createFlowInputSchema(context.getField("inputParams").value.value)
                else
                    return null;
                
            default:
                return null;
        }
    }

    validate = (fieldName: string, context: ITriggerContribution): Observable<IValidationResult> | IValidationResult => {
        let input = context.getField("hasObservers").value;
        switch (fieldName) {
            case "observerManual":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    vresult.setVisible(input);
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
        let initrigger = modelService.createTriggerElement("Dovetail-CorDApp/R3FlowInitiator");
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
                } else {
                    let inputp = context.getField("inputParams");
                    if(Boolean(inputp.value)){
                        initrigger.handler.settings[s].value = {
                            "metadata": "",
                            "value": context.getField("inputParams").value.value
                        };
                    }
                }
            }
           
        }

        let flowName = context.getFlowName();
        let iniflowModel = modelService.createFlow(flowName, context.getFlowDescription());
        let builder = modelService.createFlowElement("Dovetail-CorDApp/txnbuilder");
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