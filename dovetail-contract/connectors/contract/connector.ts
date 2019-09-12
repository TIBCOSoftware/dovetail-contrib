import { Injectable, Injector } from "@angular/core";
import { WiContrib, WiServiceHandlerContribution, WiContributionUtils,AUTHENTICATION_TYPE } from "wi-studio/app/contrib/wi-contrib";
import { IConnectorContribution } from "wi-studio/common/models/contrib";
import { IValidationResult, ValidationResult, ValidationError } from "wi-studio/common/models/validation";
import { IActionResult, ActionResult } from "wi-studio/common/models/contrib";
import { Observable } from "rxjs/Observable";
import {Http} from "@angular/http";
const zstring = require("./lz-string");

@Injectable()
@WiContrib({})
export class ContractConnectorService extends WiServiceHandlerContribution {
    constructor( private injector: Injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IConnectorContribution): Observable<any> | any => {
       
        return null;    
    }

    validate = (fieldName: string, context: IConnectorContribution): Observable<IValidationResult> | IValidationResult => {
       
        return null;
    }

    action = (name: string, context: IConnectorContribution): Observable<IActionResult> | IActionResult => {

       if(name === "Done"){
            var contracts = this.loadContract(context)
            if(contracts){
               // var json = JSON.parse(contracts)
                var txns = []
                var schemas = []
                for (const c of contracts){
                    txns.push(c.transaction)
                    
                    schemas.push([c.transaction, zstring.compressToUTF16(c.schema)])
                }
                for(let i=0; i<context.settings.length; i++){
                    if(context.settings[i].name==="transactions"){
                        context.settings[i].value = txns;
                    } else if(context.settings[i].name==="schemas"){
                        context.settings[i].value = schemas;
                    } 
                }
            }
      
            return Observable.create(observer => {
                let actionResult = {
                    context: context,
                    authType: AUTHENTICATION_TYPE.BASIC,
                    authData: {}
                }
                observer.next(ActionResult.newActionResult().setResult(actionResult));
            })
        } 
        return null
          
    }

    getFieldValue(context, fieldName): any {
        for(let i=0; i<context.settings.length; i++){
            if(context.settings[i].name=== fieldName){
                return context.settings[i].value;
            }
        }
    }

    loadContract(context){
        var model = context.getField("contractFile").value
        if(model){
            var json = JSON.parse(atob(model.content.split(",")[1].trim()))
            var contracts = []

            for(var t of json.triggers) {
                if(t.ref === "#action"){
                    for(var h of t.handlers){
                        var schema = JSON.parse(h.schemas.output.transactionInput.value)
                        var meta = JSON.parse(schema.description)
                        var ns = meta.metadata.asset.substring(0, meta.metadata.asset.lastIndexOf("."))
                        var contract = {transaction: ns + "." + h.action.settings.flowURI.substring(11), contract: meta.metadata.asset+"Contract", contractState: meta.metadata.asset, schema:h.schemas.output.transactionInput.value}
                        contracts.push(contract)
                    }
                }
            }
            return contracts
            //return JSON.stringify(contracts, null, 2)
        }
        return null
    }
  
}
