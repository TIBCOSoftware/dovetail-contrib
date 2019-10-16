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
export class ContractTriggerHandler extends WiServiceHandlerContribution {
    
    constructor(private injector: Injector, private http: Http, private contribModelService: WiContribModelService) {
        super(injector, http, contribModelService);
    }
    
    value = (fieldName: string, context: ITriggerContribution): Observable<any> | any => {
        let conId = context.getField("asset").value;
   
        switch(fieldName) {
            case "asset":
                let connectionRefs = [];
                
                return Observable.create(observer => {
                    WiContributionUtils.getConnections(this.http, "Dovetail-Ledger").subscribe((data: IConnectorContribution[]) => {
                        data.forEach(connection => {
                            if ((<any>connection).isValid) {
                                if(connection.name === "AssetSchemaConnector"){
                                    connectionRefs.push({
                                        "unique_id": WiContributionUtils.getUniqueId(connection),
                                        "name": this.getSettingValue(connection, "displayname")
                                    });
                                }
                            }
                        });
                        observer.next(connectionRefs);
                    });
                });
                
            case "actors":
                var actors = context.getField("actors").value
                if(Boolean(conId) == false || Boolean(actors) == false || Boolean(actors.value))
                    return null;

                return Observable.create(observer => {
                    WiContributionUtils.getConnection(this.http, conId)
                                        .map(data => data)
                                        .subscribe(data => {
                                            var party = []
                                            var json1 = JSON.parse(this.getSettingValue(data, "schema"))
                                            var json2 = JSON.parse(json1.description)
                                            var available = Array.from(new Set(json2.metadata.issueSigners.concat(json2.metadata.exitSigners).concat(json2.metadata.participants)))
                                            
                                            for (var p of available) {
                                                party.push({party:p, certAttributes:""})
                                            }
                                            
                                           if(party.length > 0)
                                                observer.next(JSON.stringify(party));
                                            else
                                                observer.next(null)
                                        });
                });
           
            case "assets":
            case "assets1":
                return Observable.create(observer => {
                    this.getAllAssets().subscribe((schemas: Map<string, string>) => {
                        var assets = []
                        for (let key of Array.from(schemas.keys())) {
                            assets.push(key)
                        }
                            
                        observer.next(JSON.stringify(assets, null, 2));
                        
                    });              
                });
            case "assetname":
                if(Boolean(conId) == false)
                    return null;

                return Observable.create(observer => {
                        WiContributionUtils.getConnection(this.http, conId)
                                            .map(data => data)
                                            .subscribe(data => {
                                                var name = this.getSettingValue(data, "name")
                                                var module = this.getSettingValue(data, "module")
                                                var fields = [{name: this.getSettingValue(data, "name"), type: "AssetRef", assetName: module + "." + name, consuming:"False", repeating:"False"}]
                                                observer.next(module + "." + name);
                                            });
                    });
                case "transactionInput":
                    if(Boolean(context.getField("input").value)== false)
                        return null
                       
                    return Observable.create(observer => {
                        this.getAllAssets().subscribe((schemas: Map<string, string>) => {
                            var txnschema = this.createSchema(context, "input",schemas)
                            observer.next(txnschema);
                        });              
                    });
                case "data":
                    if(Boolean(context.getField("reply").value)== false)
                        return null
                        
                    return Observable.create(observer => {
                        this.getAllAssets().subscribe((schemas: Map<string, string>) => {
                            var txnschema = this.createSchema(context, "reply", schemas)
                            observer.next(txnschema);
                        });              
                    });
                case "assetschema":
                    if(Boolean(conId) == false)
                        return null

                    return Observable.create(observer => {
                        WiContributionUtils.getConnection(this.http, conId)
                                            .map(data => data)
                                            .subscribe(data => {
                                                var schema = this.getSettingValue(data, "schema")
                                                observer.next(schema);
                                            });
                    });
                case "from":
                case "until":
                    var timewindow = context.getField("timewindow").value
                    if(timewindow !== "Any time"){
                        if(Boolean(conId) == false)
                            return null

                        return Observable.create(observer => {
                            WiContributionUtils.getConnection(this.http, conId)
                                            .map(data => data)
                                            .subscribe(data => {
                                                var fields = []
                                                var schema = this.getSettingValue(data, "schema")
                                                fields = this.getDateTimeFields(null, JSON.parse(schema).properties, fields)
                                            /*    var meta = JSON.parse(JSON.parse(schema).description)
                                                for (var a of meta.attributes){
                                                    if(a.type === "com.tibco.dovetail.system.Instant")
                                                        fields.push(a.name)
                                                }*/
                                                observer.next(fields);
                                            });
                            });
                    }
                    
            default: 
                return null;
        }
            
    }

    validate = (fieldName: string, context: ITriggerContribution): Observable<IValidationResult> | IValidationResult => {
        var timewindow = context.getField("timewindow").value
        switch(fieldName){
            case "from":      
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    vresult.setVisible(timewindow === "Only valid if after..." || timewindow === "Only valid if between...");
                    observer.next(vresult);
                });
            case "until":      
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    vresult.setVisible(timewindow === "Only valid if before..." || timewindow === "Only valid if between...");
                    observer.next(vresult);
                });
            case "assets1":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    vresult.setVisible(context.getMode() === MODE.WIZARD || context.getMode() === MODE.SERVERLESS_FLOW);
                    observer.next(vresult);
                });
        }
        return null;
    }

    action = (actionId: string, context: ICreateFlowActionContext): Observable<IActionResult> | IActionResult => {
       
        let result = CreateFlowActionResult.newActionResult();
        return Observable.create(observer => {
            this.createFlow(context, result);
            let actionResult = ActionResult.newActionResult().setSuccess(true).setResult(result);
            observer.next(actionResult);      
        });
    }

    createFlow(context, result) : string{
        let modelService = this.getModelService();
        let trigger = modelService.createTriggerElement("Dovetail-Ledger/ContractActionTrigger");
        if (trigger) {
            for(let t = 0; t < trigger.handler.settings.length; t++) {
                if (trigger.handler.settings[t].name === "asset" ) {
                    trigger.handler.settings[t].value = context.getField("asset").value;
                } else if (trigger.handler.settings[t].name === "txnType") {
                    trigger.handler.settings[t].value = context.getField("txnType").value;
                } else if (trigger.handler.settings[t].name === "actors") {
                    trigger.handler.settings[t].value = context.getField("actors").value;
                } else if (trigger.handler.settings[t].name === "input") {
                    trigger.handler.settings[t].value = context.getField("input").value;
                } else if (trigger.handler.settings[t].name === "assetname") {
                    trigger.handler.settings[t].value = context.getField("assetname").value;
                } else if (trigger.handler.settings[t].name === "assetschemas") {
                    trigger.handler.settings[t].value = context.getField("assetschemas").value;
                } else if (trigger.handler.settings[t].name === "timewindow") {
                    trigger.handler.settings[t].value = context.getField("timewindow").value;
                } else if (trigger.handler.settings[t].name === "from") {
                    trigger.handler.settings[t].value = context.getField("from").value;
                } else if (trigger.handler.settings[t].name === "until") {
                    trigger.handler.settings[t].value = context.getField("until").value;
                } else if (trigger.handler.settings[t].name === "reply") {
                    trigger.handler.settings[t].value = context.getField("reply").value;
                } 
            }
            for (let j = 0; j < trigger.outputs.length; j++) {
                if (trigger.outputs[j].name === "transactionInput") {
                    trigger.outputs[j].value =  context.getField("transactionInput").value
                }
            }

            for (let j = 0; j < trigger.reply.length; j++) {
                if (trigger.reply[j].name === "data") {
                    var reply = context.getField("data").value
                    if(reply){
                        trigger.reply[j].value =  reply
                    }
                    else{
                        trigger.reply[j].value =  {
                            "value":{}
                        }
                    }
                    break
                }
            }
        }

        let flowName = context.getFlowName();
        let flowModel = modelService.createFlow(flowName, context.getFlowDescription());
       
        result = result.addTriggerFlowMapping(lodash.cloneDeep(trigger), lodash.cloneDeep(flowModel));
        return flowName;
    }
    
    getSettingValue(connection, setting):string {
        for(let i=0; i < connection.settings.length; i++) {
            if(connection.settings[i].name === setting){
                return connection.settings[i].value
            }
        }
    }

    getAllAssets():Observable<Map<string, string>>  {
        var schemas = new Map()
        return Observable.create(observer => {
            WiContributionUtils.getConnections(this.http, "Dovetail-Ledger").subscribe((data: IConnectorContribution[]) => {
                data.forEach(connection => {
                    if ((<any>connection).isValid) {
                        if(connection.name === "AssetSchemaConnector"){
                            var name = this.getSettingValue(connection, "name")
                            var module = this.getSettingValue(connection, "module")
                            var schema = this.getSettingValue(connection, "schema")
                          schemas.set(module+"."+name, schema)
                        }
                    }
                });
                observer.next(schemas);
            });
        });
    }

    getAllStructs():Observable<Map<string, string>>  {
        var schemas = new Map()
        return Observable.create(observer => {
            WiContributionUtils.getConnections(this.http, "Dovetail-General").subscribe((data: IConnectorContribution[]) => {
                data.forEach(connection => {
                    if ((<any>connection).isValid) {
                        if(connection.name === "StructSchema"){
                            var name = this.getSettingValue(connection, "name")
                            var module = this.getSettingValue(connection, "module")
                            var schema = this.getSettingValue(connection, "schema")
                          schemas.set(module+"."+name, schema)
                        }
                    }
                });
                observer.next(schemas);
            });
        });
    }

    createSchema(context, field, schemas) : string{
        var input = context.getField(field).value.value
        if(input){
            var fields = JSON.parse(input);
            var schema = {schema:"http://json-schema.org/draft-04/schema#", type: "object", properties:{}}
            var metadata = {metadata: {type:"Transaction", parent:"", actors:[], asset:"", timewindow:{}}, attributes:[]}
            metadata.metadata.asset = context.getField("assetname").value
            if(field === "input"){
                var authorized = context.getField("actors").value.value
                if(authorized){
                    var authjson = JSON.parse(authorized)
                    for(var i=0; i<authjson.length; i++){
                        metadata.metadata.actors.push(authjson[i].party + "|" + authjson[i].certAttributes)
                    }
                }
                var timewindow = context.getField("timewindow").value
                if(timewindow === "Only valid if after..." || timewindow === "Only valid if between...")
                    metadata.metadata.timewindow["from"] = context.getField("from").value

                if(timewindow === "Only valid if before..." || timewindow === "Only valid if between...")
                    metadata.metadata.timewindow["from"] = context.getField("until").value
            }

            for(var i=0; i<fields.length; i++){
                    let name = fields[i].name;
                    let tp = fields[i].type;
                    let repeating = fields[i].repeating;
                    let isArray = false;
                    let isRef = false;
                    var isReferenceData = false
                    var isAsset = false
                    var isParticipant = false
                    let attr = {};
                    let datatype = {type: tp.toLowerCase()};
                    let systype = tp;
                    switch (tp) {
                        case "Party":
                            datatype.type = "string";
                            systype = "com.tibco.dovetail.system.Party";
                            isRef = true;
                            isParticipant = true
                            break;
                        case "LinearId":
                            datatype.type = "string";
                            systype = "com.tibco.dovetail.system.UniqueIdentifier";
                            break;
                        case "Amount<Currency>":
                            datatype.type = "object";
                            datatype["properties"] = {currency: {type: "string"}, quantity: {type: "number"}};
                            systype = "com.tibco.dovetail.system.Amount<Currency>";
                            break;
                        case "Integer":
                        case "Long":
                            datatype.type = "number";
                            break;
                        case "Decimal":
                            datatype.type = "string"
                            systype = "Decimal"
                            break;
                        case "LocalDate":
                            datatype.type = "string";
                            datatype["format"] = "date-time";
                            systype = "com.tibco.dovetail.system.LocalDate"
                            break;
                        case "DateTime":
                            datatype.type = "string";
                            datatype["format"] = "date-time";
                            systype = "com.tibco.dovetail.system.Instant"
                            break;
                        case "AssetRef":
                        case "Record":
                            if(tp == "AssetRef")
                                isRef = true
                            systype = fields[i].recordType
                            datatype.type = "object"
                            var asset = schemas.get(systype)
                            if(asset){
                                var aschema = JSON.parse(schemas.get(systype))
                                datatype["properties"] = aschema.properties
                            } else {
                                console.log("error: cann't find asset - " + systype)
                            }
                            if(field === "input" && fields[i].consuming == "False")
                                isReferenceData = true

                            isAsset = true
                            break
                        case "AssetRef<Cash>":
                            systype = "com.tibco.dovetail.system.Cash"
                            isRef = true
                            isAsset = true
                            datatype.type = "object"
                            datatype["properties"] = {amt:{type:"object", properties:{currency: {type: "string"}, quantity: {type: "number"}}}, issuer:{type:"string"}, issuerRef:{type:"string"}, owner:{type:"string"}}
                    }
                    if(repeating === "True"){
                        schema.properties[name] = {type: "array", items: datatype}
                        isArray = true;
                    } else {
                        schema.properties[name] = datatype
                    }
        
                    attr["name"] = name;
                    attr["type"] = systype;
                    attr["isRef"] = isRef
                    attr["isArray"] = isArray;
                    attr["isReferenceData"] = isReferenceData
                    attr["isAsset"] = isAsset
                    attr["isParticipant"] = isParticipant
                    metadata.attributes.push(attr);
            }
            schema["description"] = JSON.stringify(metadata);
            return JSON.stringify(schema);
        } else {
            return null;
        }
    }

    getRequiredAssetSchemas(context, schemas) {
        var input = context.getField("input").value.value
        var assetschemas = []
        if(input){
            var fields = JSON.parse(input);
           
            for(var i=0; i<fields.length; i++){
                    
                let tp = fields[i].type;
                
                switch (tp) {
                    case "AssetRef":
                    case "Asset":
                        var systype = fields[i].recordType
                        assetschemas.push([systype, schemas.get(systype)])
                        break
                }
            } 
        } 
        return assetschemas;
    }

    getDateTimeFields(parent, properties, fields){
       var keys = Object.keys(properties)
        for(var k of keys) {
            var p = properties[k]
            console.log(p)
            switch(p.type){
                case "string":
                    if(typeof(p["format"]) !== 'undefined' && p.format === "date-time"){
                        if(parent)
                            fields.push(parent + "." + k)
                        else
                            fields.push(k)
                    }
                    break
                case "object":
                    if(parent)
                        fields = this.getDateTimeFields(parent + "." + k, p.properties, fields)
                    else
                        fields = this.getDateTimeFields(k, p.properties, fields)
                    
                    break
                case "array":
                    if(parent)
                        fields = this.getDateTimeFields(parent + "." + k, p.items.properties, fields)
                    else
                        fields = this.getDateTimeFields(k, p.items.properties, fields)
                   
                    break
            }
            
        }
        return fields
    }
}