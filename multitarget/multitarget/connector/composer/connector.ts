/// <reference path="./bundle.d.ts" />
/// <amd-dependency path="./bundle"/>
import { Inject, Injectable, Injector } from "@angular/core";
import { WiContrib, WiServiceHandlerContribution, WiContributionUtils,AUTHENTICATION_TYPE } from "wi-studio/app/contrib/wi-contrib";
import { IConnectorContribution } from "wi-studio/common/models/contrib";
import { IValidationResult, ValidationResult, ValidationError } from "wi-studio/common/models/validation";
import { IActionResult, ActionResult } from "wi-studio/common/models/contrib";
import { Observable } from "rxjs/Observable";
const modelParser = require("./bundle");

@Injectable()
@WiContrib({})
export class ComposerConnectorService extends WiServiceHandlerContribution {
    constructor( @Inject(Injector) injector) {
        super(injector);
    }

    value = (fieldName: string, context: IConnectorContribution): Observable<any> | any => {
        return null;    
    }

    validate = (fieldName: string, context: IConnectorContribution): Observable<IValidationResult> | IValidationResult => {
        let loadBna = this.getFieldValue(context, "mode");
        switch(fieldName) {
            case "enterModel":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    vresult.setVisible(!loadBna);
                    let modelV = this.getFieldValue(context, "enterModel");
                    if(!loadBna && (modelV == false || modelV.value == "") ){
                        vresult.setValid(false).setError("REQUIRED_VALUE_NOT_SET", "Please enter the composer model");
                    }
                    observer.next(vresult);
                });
            case "modelFile":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    vresult.setVisible(loadBna);
                    let modelV = this.getFieldValue(context, "modelFile");
                    if(loadBna && (modelV == false || modelV == "") ){
                        vresult.setValid(false).setError("REQUIRED_VALUE_NOT_SET", "Please select composer model archive file");
                    }
                    observer.next(vresult);
                });
            case "assets":
            case "transactions":
            case "schemas":
            case "events":
            case "participants":
            case "concepts":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    observer.next(vresult.setVisible(false));
                });
            default:
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    observer.next(vresult);
                });
        }
       
    }

    action = (name: string, context: IConnectorContribution): Observable<IActionResult> | IActionResult => {
        let assets = [];
        let concepts = [];
        let txns = [];
        let events = [];
        let participants = [];
        let bna;
        let schemas = new Map();
       if(name === "Save Model"){
            return Observable.create(observer => {
                try {
                    if(this.getFieldValue(context, "mode")) {
                        bna = this.getFieldValue(context, "modelFile");
                        modelParser.fromArchive(bna.content.split(",")[1].trim(), assets, txns, schemas, events, participants, concepts).then(() => {
                            observer.next(this.processModelOutput(context,assets, txns, schemas, events,participants, concepts));
                        });
                    } else {
                        modelParser.fromText(this.getFieldValue(context, "enterModel").value, assets, txns, schemas, events, participants, concepts);
                        observer.next(this.processModelOutput(context,assets, txns, schemas, events, participants, concepts));
                    }
                } catch(err) {
                    console.log(name + ":" + err.message);
                    return Observable.create(observer => {
                        observer.next(ActionResult.newActionResult().setSuccess(false).setResult(new ValidationError("chaincode-1002", "Action failed :" + err.message)));
                        observer.complete();
                    });
                }
            });
        } else {
            return Observable.create(observer => {
                observer.next(ActionResult.newActionResult().setSuccess(false).setResult(new ValidationError("chaincode-1001", "Invalid action :" + name)));
            });
        }
    }

    processModelOutput(context, assets,txns, schemas, events, participants, concepts) : IActionResult{
        for(let i=0; i<context.settings.length; i++){

            if(context.settings[i].name==="assets"){
                context.settings[i].value = assets;
            } else if(context.settings[i].name==="transactions"){
                context.settings[i].value = txns;
            } else if(context.settings[i].name==="schemas"){
                context.settings[i].value = schemas;
            } else if(context.settings[i].name==="events"){
                context.settings[i].value = events;
            } else if(context.settings[i].name==="participants"){
                context.settings[i].value = participants;
            } else if(context.settings[i].name==="concepts"){
                context.settings[i].value = concepts;
            }
        }
        
        let actionResult = {
            context: context,
            authType: AUTHENTICATION_TYPE.BASIC,
            authData: {}
        }
        return ActionResult.newActionResult().setResult(actionResult);
    }

    getFieldValue(context, fieldName): any {
        for(let i=0; i<context.settings.length; i++){
            if(context.settings[i].name=== fieldName){
                return context.settings[i].value;
            }
        }
    }

    setStatus(context, status){
        for(let i=0; i<context.settings.length; i++){
            if(context.settings[i].name=== "status"){
                context.settings[i].value = status;
            }
        }
    }
}
