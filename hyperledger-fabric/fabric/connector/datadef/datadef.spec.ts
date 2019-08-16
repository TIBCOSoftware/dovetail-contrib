
import { Injector } from "@angular/core";
import { HttpModule, Http, BaseRequestOptions, XHRBackend, Response, ResponseOptions } from "@angular/http";
import { } from "jasmine";
import { TestBed, inject, fakeAsync, tick } from "@angular/core/testing";
import { MockBackend, MockConnection } from "@angular/http/testing";
import { datadefHandler } from "./datadefHandler";
import { IContributionTypes, WiServiceContribution, IFieldDefinition, WiContributionUtils } from "wi-studio/index";
import * as TypeMoq from "typemoq";

/**
 * datadefHandler tests
 */
export let t1 = describe("datadefHandler tests", () => {

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpModule],
      providers: [
        { provide: WiServiceContribution, useClass: datadefHandler },
        { provide: XHRBackend, useClass: MockBackend }
      ]
    });
  });


  /**
   * Test datadefHandler
   */
  describe("datadefHandler", () => {
    it("should return datadefHandler", () => {
      inject([Injector, Http], (injector: Injector, http: Http) => {
        let svc = new datadefHandler(injector);
        expect(svc !== null).toBeTruthy("datadefHandler not found");
      })();
    });
  });

});
