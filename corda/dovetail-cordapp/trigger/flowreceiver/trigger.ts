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

@WiContrib({})
@Injectable()
export class R3FlowReceiverTriggerHandler extends WiServiceHandlerContribution {
    
    constructor(private injector: Injector, private http: Http, private contribModelService: WiContribModelService) {
        super(injector, http, contribModelService);
    }
    
    value = (fieldName: string, context: ITriggerContribution): Observable<any> | any => {
        return null;
            
    }

    validate = (fieldName: string, context: ITriggerContribution): Observable<IValidationResult> | IValidationResult => {
        let input = context.getField("flowType").value;
        switch (fieldName) {
            case "useAnonymousIdentity":
            case "output":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    if(input==="receiver")
                        vresult.setVisible(true);
                    else
                        vresult.setVisible(false);

                    observer.next(vresult);
                });
        }
        return null;
    }

    action = (actionId: string, context: ICreateFlowActionContext): Observable<IActionResult> | IActionResult => {
       
        let result = CreateFlowActionResult.newActionResult();
        let flows = []
        return Observable.create(observer => {
            let flowName = context.getFlowName();
            this.createFlow(context, flowName, result);                                     
            let actionResult = ActionResult.newActionResult().setSuccess(true).setResult(result);
            observer.next(actionResult);
        });
    }

    createFlow(context, flowName, result) : string{
        let modelService = this.getModelService();
        let trigger = modelService.createTriggerElement("Dovetail-CorDApp/R3FlowReceiver");
        if (trigger) {
            for (let s = 0; s < trigger.handler.settings.length; s++) {
                if (trigger.handler.settings[s].name === "flowType") {
                    trigger.handler.settings[s].value = context.getField("flowType").value;
                } else if (trigger.handler.settings[s].name === "initiatorFlow") {
                    trigger.handler.settings[s].value = context.getField("initiatorFlow").value;
                } else if(trigger.handler.settings[s].name === "useAnonymousIdentity"){
                    trigger.handler.settings[s].value = context.getField("useAnonymousIdentity").value;
                }
            }
        }
        let flowModel = modelService.createFlow(flowName, context.getFlowDescription());
        result = result.addTriggerFlowMapping(lodash.cloneDeep(trigger), lodash.cloneDeep(flowModel));
        return flowName;
    }
}