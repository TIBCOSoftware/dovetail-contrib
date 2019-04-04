import {Observable} from "rxjs/Observable";
import {Injectable, Injector, Inject} from "@angular/core";
import {Http} from "@angular/http";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    IFieldDefinition,
    IActivityContribution,
    ActionResult,
    IActionResult,
    WiContributionUtils,
    WiContribModelService,
    IConnectorContribution
} from "wi-studio/app/contrib/wi-contrib";


@WiContrib({})
@Injectable()
export class AnalyticsActivityContributionHandler extends WiServiceHandlerContribution {
  /*  constructor(private injector: Injector, private http: Http,) {
        super(injector, http);
    }
   */

  constructor(private injector: Injector, private http: Http) {
    super(injector, http);
}
       
    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {
        let action = context.getField("operation").value;
        let asset = context.getField("assetName").value;
        let isArray = false;
        
        let conId = context.getField("model").value;
    
        switch(fieldName){
            case "input":
                if(Boolean(conId) == false || Boolean(asset) == false)
                    return null;

                return Observable.create(observer => {
                    this.getSchemas(conId).subscribe( schemas => {
                        observer.next(schemas[asset]);
                        
                    });
                });
            case "output":
                if(Boolean(conId) == false || Boolean(asset) == false)
                    return null;

                return Observable.create(observer => {
                    this.getSchemas(conId).subscribe( schemas => {
                        observer.next(schemas[asset]);
                    });
                });
            default:
                return null;
        }
    }
    getSchemas(conId):  Observable<any> {
        let schemas = new Map();
        return Observable.create(observer => {
            WiContributionUtils.getConnection(this.http, conId)
                            .map(data => data)
                            .subscribe(data => {
                                let schemas = new Map();
                                for (let setting of data.settings) {
                                    if(setting.name === "schemas") {
                                        setting.value.map(item => schemas[item[0]] = item[1]);
                                        observer.next(schemas);
                                        break;
                                    }
                                }
                            });
                        });
    }

}