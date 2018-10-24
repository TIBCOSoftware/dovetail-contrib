import { Http } from "@angular/http";
import {WiContrib, WiServiceProviderContribution, AbstractContribFieldProvider, AbstractContribValidationProvider} from "wi-studio/app/contrib/wi-contrib";
import {IActivityContribution, IContributionContext} from "wi-studio/common/models/contrib";
import {IValidationResult, ValidationResult} from "wi-studio/common/models/validation";
import { Injectable, Inject, Injector } from "@angular/core";
import {Observable} from "rxjs/Observable";

@Injectable()
export class Field1Provider extends AbstractContribFieldProvider {

    constructor( @Inject(Injector) private injector, @Inject(Http) private http) {
        super();
    }
    getFieldValue(context: IActivityContribution): Observable<string[]> {
        return Observable.create(observer => {
            try {
                observer.next([context.name]);
            } finally {
                observer.complete();
            }
        });
    }
    getHttp() { return this.http; }
    getInjector() { return this.injector; }
}
@Injectable()
export class TestInjectable {

}
@Injectable()
export class Field1ValidationProvider extends AbstractContribValidationProvider {

    private platformRef = null;
    constructor() {
        super();
    }

    validate(context: IContributionContext): Observable<IValidationResult> {
        return Observable.create(observer => {
            try {
                // tslint:disable-next-line:no-bitwise
                let vresult: IValidationResult = ValidationResult.newValidationResult();
                observer.next(vresult);
            } finally {
                observer.complete();
            }
        });
    }
    getPlatformRef() { return this.platformRef; }
}

@WiContrib({
    validationProviders: [
        {
            field: "field1",
            useClass: Field1ValidationProvider
        }
    ],
    fieldProviders: [
        {
            field: "field1",
            useClass: Field1Provider
        }
    ]
})
@Injectable()
export class KVLedgerActivityService extends WiServiceProviderContribution {
    constructor( @Inject(Injector) injector, @Inject(Http) http: Http) {
        super(injector, http);
    }
}

