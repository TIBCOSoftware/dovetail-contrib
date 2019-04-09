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

}