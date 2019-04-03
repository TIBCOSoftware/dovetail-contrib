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
export class R3FlowInitiatorTriggerHandler extends WiServiceHandlerContribution {
    
    constructor(private injector: Injector, private http: Http, private contribModelService: WiContribModelService) {
        super(injector, http, contribModelService);
    }
    
    value = (fieldName: string, context: ITriggerContribution): Observable<any> | any => {
        switch(fieldName) {
            case "transactionInput":
               let params = context.getField("inputParams").value;
               return this.createFlowInputSchema(params.value);
            default:
                return null;
        }
       
            
    }

    validate = (fieldName: string, context: ITriggerContribution): Observable<IValidationResult> | IValidationResult => {
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
                if (initrigger.handler.settings[s].name === "allowRPCClient") {
                    initrigger.handler.settings[s].value = context.getField("allowRPCClient").value;
                } else if (initrigger.handler.settings[s].name === "useAnonymousIdentity") {
                    initrigger.handler.settings[s].value = context.getField("useAnonymousIdentity").value;
                } else {
                    let inputp = context.getField("inputParams");
                    initrigger.handler.settings[s].value = {
                        "metadata": "",
                        "value": context.getField("inputParams").value
                    };
                }
            }
            for (let j = 0; j < initrigger.outputs.length; j++) {
                if (initrigger.outputs[j].name === "transactionInput") {
                    initrigger.outputs[j].value =  {
                        "value": this.createFlowInputSchema(context.getField("inputParams").value),
                        "metadata": ""
                    };
                    break;
                }
            }
        }

        let flowName = context.getFlowName();
        let iniflowModel = modelService.createFlow(flowName+"Initiator", context.getFlowDescription());
        let builder = modelService.createFlowElement("CorDApp/txnbuilder");
        iniflowModel.addFlowElement(builder);
        let sign = modelService.createFlowElement("CorDApp/signandcommit");
        iniflowModel.addFlowElement(sign);
        result = result.addTriggerFlowMapping(lodash.cloneDeep(initrigger), lodash.cloneDeep(iniflowModel));

        let rectrigger = modelService.createTriggerElement("CorDApp/R3FlowReceiver");
        if (rectrigger) {
            for (let j = 0; j < rectrigger.settings.length; j++) {
                if (rectrigger.handler.settings[j].name === "initiatorFlow") {
                    rectrigger.handler.settings[j].value = flowName + "Initiator"; 
                    break;
                }
            }
        }
        
        let recflowModel = modelService.createFlow(flowName+"Responder", context.getFlowDescription());
        let recsign = modelService.createFlowElement("CorDApp/receiversign");
        recflowModel.addFlowElement(recsign);
        result = result.addTriggerFlowMapping(lodash.cloneDeep(rectrigger), lodash.cloneDeep(recflowModel));
        return flowName;
    }

    createFlowInputSchema(inputParams):String {
        if(Boolean(inputParams) == false)
            return "{}";

       let inputs = JSON.parse(inputParams);
       let schema = {schema:"http://json-schema.org/draft-04/schema#", type: "object", properties:{}}
       let metadata = {metadata: {type: "Transaction"}, attributes: []};

       if(inputs) {
           for(let i=0; i<inputs.length; i++){
                let name = inputs[i].parameterName;
                let tp = inputs[i].type;
                let repeating = inputs[i].repeating;
                let datatype = {type: tp.toLowerCase()};
                let javatype = tp;
                let isRef = false;
                let isArray = false;
                let attr = {};

                switch (tp) {
                    case "Party":
                        datatype.type = "string";
                        javatype = "net.corda.core.identity.Party";
                        isRef = true;
                        break;
                    case "LinearId":
                        datatype.type = "string";
                       // datatype.type = "object";
                      //  datatype["properties"] = {uuid: {type: "string"}, externalId: {type: "string"}};
                        javatype = "net.corda.core.contracts.UniqueIdentifier";
                        break;
                    case "Amount<Currency>":
                        datatype.type = "object";
                        datatype["properties"] = {currency: {type: "string"}, quantity: {type: "number"}};
                        javatype = "net.corda.core.contracts.Amount<Currency>";
                        break;
                    case "Integer":
                    case "Long":
                        datatype.type = "number";
                        break;
                }
                if(repeating === "true"){
                    schema.properties[name] = {type: "array", items: {datatype}}
                    isArray = true;
                } else {
                    schema.properties[name] = datatype
                }

                attr["name"] = name;
                attr["type"] = javatype;
                attr["isRef"] = isRef
                attr["isArray"] = isArray;
                metadata.attributes.push(attr);
           }
           schema["description"] = JSON.stringify(metadata);
           return JSON.stringify(schema);
       } else {
           return "{}";
       }
    }
}